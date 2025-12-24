package c2

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/hardal7/pex/internal/config"
	logger "github.com/hardal7/pex/internal/util"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "help",
	Short: "Adversary Emulation Framework",
	Long:  `pex: Open source adversary emulation framework - Made with <3.`,
}

const uuidLength int = 36

var taskCmd = &cobra.Command{
	Use:     "task",
	Short:   "Task a shell command",
	Long:    `Create tasks for selected agents to run a given shell command`,
	Example: `Example: task ls -la`,
	Run: func(cmd *cobra.Command, args []string) {
		if State.SelectedAgent.Alias == "ALL" {
			for _, agent := range State.RegisteredAgents {
				State.Tasks = append(State.Tasks, Task{Command: strings.Join(args, " "), Recipient: agent})
			}
		} else if State.SelectedAgent.Alias != "NONE" {
			State.Tasks = append(State.Tasks, Task{Command: strings.Join(args, ""), Recipient: Agent{UUID: State.SelectedAgent.UUID}})
		}
	},
}

var pickCmd = &cobra.Command{
	Use:   "pick <agent>",
	Short: "Pick an agent",
	Long: `Pick an agent to submit the tasks to
Valid identifiers are: UUID, Hostname, Alias`,
	Example: `Example: select 192.168.1.2`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		selectName := args[0]
		if len(selectName) == uuidLength {
			State.SelectedAgent.UUID = selectName
		} else if strings.Contains(selectName, ":") || strings.Contains(selectName, ".") {
			State.SelectedAgent.Hostname = selectName
		} else {
			State.SelectedAgent.Alias = selectName
		}
		if selectName == "ALL" {
			logger.Info("Selected all agents")
		} else {
			for _, agent := range State.RegisteredAgents {
				if State.SelectedAgent.UUID == agent.UUID || State.SelectedAgent.Hostname == agent.Hostname || State.SelectedAgent.Alias == agent.Alias {
					State.SelectedAgent = agent
					logger.Info("Selected agent with " + logAgent(State.SelectedAgent))
					break
				} else {
					logger.Info("No registered agent found with UUID: " + State.SelectedAgent.UUID)
					State.SelectedAgent.Alias = "NONE"
				}
			}
		}
	},
}

var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "List agents",
	Long:  `List agents registered to the server`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(State.RegisteredAgents) == 0 {
			logger.Info("No agents registered")
		} else {
			for i := range len(State.RegisteredAgents) {
				logger.Info("Agent " + strconv.Itoa(i) + ": " + logAgent(State.RegisteredAgents[i]))
			}
		}
	},
}

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List queued tasks",
	Long:  `List queued tasks waiting to be run by agents`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(State.Tasks) == 0 {
			logger.Info("No tasks in queue")
		} else {
			for i, task := range State.Tasks {
				logger.Info("Task " + strconv.Itoa(i) + ": " + " Command: " + task.Command + " Recipient: " + task.Recipient.UUID)
			}
		}
	},
}

var aliasCmd = &cobra.Command{
	Use:   "alias <identifier> <alias>",
	Short: "Give an agent alias",
	Long: `Give an agent with specified identifier an alias
Valid identifiers are: UUID, Hostname, Username`,
	Example: `Example: alias 192.168.1.1 host1`,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		for i, agent := range State.RegisteredAgents {
			if agent.UUID == args[0] || agent.Hostname == args[0] || agent.Username == args[0] {
				State.RegisteredAgents[i].Alias = args[1]
				logger.Info("Aliased agent to name: " + State.RegisteredAgents[i].Alias)
				break
			}
			logger.Info("No agents found with given identifier")
		}
	},
}

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Start a session",
	Long:  `Signal an agent to start a TCP session`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Initiating session")
		// TODO: Switch between connection types
		go InitiateSession()
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a teamserver",
	Long:  `Start a teamserver to host the server to multiple clients`,
	Run: func(cmd *cobra.Command, args []string) {
		if State.IsServing != true {
			go HostTeamserver()
		} else {
			logger.Info("Teamserver is already running on port: " + config.TeamserverPort)
		}
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate <target> <name>",
	Short: "Generate a beacon",
	Long: `Generate a beacon with given name to given target architecture
Valid targets are: windows, linux, macos`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// FIXME: Only works if inside root directory
		build := exec.Command("go", "build", "-o", args[1], "./cmd/agent/main.go")
		build.Env = os.Environ()
		build.Env = append(build.Env, "GOARCH=amd64")
		switch args[0] {
		case "windows":
			build.Env = append(build.Env, "GOOS=windows")
		case "linux":
			build.Env = append(build.Env, "GOOS=linux")
		case "macos":
			build.Env = append(build.Env, "GOOS=darwin")
		default:
			logger.Info("Invalid target architecture, valid targets are: windows, linux, macos")
			return
		}
		logger.Info("Generated beacon for target: " + args[0])
	},
}
var loglevelCmd = &cobra.Command{
	Use:   "loglevel <level>",
	Short: "Change log level",
	Long: `Change the verbosity of the log level
Valid log levels are: debug, info, silent`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "debug":
			config.LogLevel = "debug"
		case "info":
			config.LogLevel = "info"
		case "silent":
			config.LogLevel = "silent"
		default:
			logger.Info("Invalid log level, valid log levels are: verbose, info, silent")
		}
		logger.Load()
	},
}

var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "Close and exit the server",
	Long:  `Close and exit the server`,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}

func logAgent(agent Agent) string {
	return ("UUID: " + agent.UUID + " Alias: " + agent.Alias + " Hostname: " + agent.Hostname + " OS: " + agent.OS)
}

func MenuCommands() *cobra.Command {
	return root
}

func FetchCommand(requestedCommand string) *cobra.Command {
	for _, command := range root.Commands() {
		if requestedCommand == strings.Split(command.Use, " ")[0] {
			return command
		}
	}
	return nil
}

func ExecuteCommand(command string, args []string) {
	fetchedCommand := FetchCommand(command)
	if fetchedCommand != nil {
		agentsCmd.Run(agentsCmd, args)
	} else {
		// TODO: Log command not found to client
	}
}

func InitCommands() {
	root.AddCommand(taskCmd, pickCmd, agentsCmd, tasksCmd, aliasCmd, sessionCmd, serveCmd, generateCmd, loglevelCmd, exitCmd)
}
