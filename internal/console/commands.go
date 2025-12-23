package console

import (
	"os"
	"strconv"
	"strings"

	"github.com/hardal7/pex/internal/c2"
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
		if c2.State.SelectedAgent.Alias == "ALL" {
			for _, agent := range c2.State.RegisteredAgents {
				c2.State.Tasks = append(c2.State.Tasks, c2.Task{Command: strings.Join(args, " "), Recipient: agent})
			}
		} else if c2.State.SelectedAgent.Alias != "NONE" {
			c2.State.Tasks = append(c2.State.Tasks, c2.Task{Command: strings.Join(args, ""), Recipient: c2.Agent{UUID: c2.State.SelectedAgent.UUID}})
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
			c2.State.SelectedAgent.UUID = selectName
		} else if strings.Contains(selectName, ":") || strings.Contains(selectName, ".") {
			c2.State.SelectedAgent.Hostname = selectName
		} else {
			c2.State.SelectedAgent.Alias = selectName
		}
		if selectName == "ALL" {
			logger.Info("Selected all agents")
		} else {
			for _, agent := range c2.State.RegisteredAgents {
				if c2.State.SelectedAgent.UUID == agent.UUID || c2.State.SelectedAgent.Hostname == agent.Hostname || c2.State.SelectedAgent.Alias == agent.Alias {
					c2.State.SelectedAgent = agent
					logger.Info("Selected agent with " + logAgent(c2.State.SelectedAgent))
					break
				} else {
					logger.Info("No registered agent found with UUID: " + c2.State.SelectedAgent.UUID)
					c2.State.SelectedAgent.Alias = "NONE"
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
		if len(c2.State.RegisteredAgents) == 0 {
			logger.Info("No agents registered")
		} else {
			for i := range len(c2.State.RegisteredAgents) {
				logger.Info("Agent " + strconv.Itoa(i) + ": " + logAgent(c2.State.RegisteredAgents[i]))
			}
		}
	},
}

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List queued tasks",
	Long:  `List queued tasks waiting to be run by agents`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(c2.State.Tasks) == 0 {
			logger.Info("No tasks in queue")
		} else {
			for i, task := range c2.State.Tasks {
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
		for i, agent := range c2.State.RegisteredAgents {
			if agent.UUID == args[0] || agent.Hostname == args[0] || agent.Username == args[0] {
				c2.State.RegisteredAgents[i].Alias = args[1]
				logger.Info("Aliased agent to name: " + c2.State.RegisteredAgents[i].Alias)
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
		go c2.InitiateSession()
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a teamserver",
	Long:  `Start a teamserver to host the server to multiple clients`,
	Run: func(cmd *cobra.Command, args []string) {
		if c2.State.IsServing != true {
			c2.State.IsServing = true
			logger.Info("Started teamserver on port: " + config.TeamserverPort)
		} else {
			logger.Info("Teamserver is already running on port: " + config.TeamserverPort)
		}
	},
}

var loglevelCmd = &cobra.Command{
	Use:   "loglevel <level>",
	Short: "Change log level",
	Long: `Change the verbosity of the log level
Valid log levels are: debug, info, silent`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		loglevel := strings.TrimSpace(strings.TrimPrefix(args[0], "loglevel "))
		switch loglevel {
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
		// TODO: Graceful exit
		os.Exit(0)
	},
}

func logAgent(agent c2.Agent) string {
	return ("UUID: " + agent.UUID + " Alias: " + agent.Alias + " Hostname: " + agent.Hostname)
}

func MenuCommands() *cobra.Command {
	return root
}

func InitCommands() {
	root.AddCommand(taskCmd, pickCmd, agentsCmd, tasksCmd, aliasCmd, sessionCmd, serveCmd, loglevelCmd, exitCmd)
}
