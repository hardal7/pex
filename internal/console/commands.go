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

func Commands() *cobra.Command {
	var root = &cobra.Command{
		Use:   "pex",
		Short: "Adversary Emulation Framework",
		Long:  `pex is an open source adversary emulation framework written in Golang for red teaming specialists and hobbyists`,
	}

	const uuidLength int = 36
	const command string = ""

	var selectCmd = &cobra.Command{
		Use:   "select",
		Short: "Select an agent",
		Long: `Usage: select <identifier>
					Valid identifiers are: UUID, Hostname, Alias
					Example: select 192.168.1.2`,
		Run: func(cmd *cobra.Command, args []string) {
			selectName := strings.TrimSpace(strings.TrimPrefix(command, "select "))
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
	root.AddCommand(selectCmd)

	var agentsCmd = &cobra.Command{
		Use:   "agents",
		Short: "List registered agents",
		Long:  `Usage: agents`,
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
	root.AddCommand(agentsCmd)

	var tasksCmd = &cobra.Command{
		Use:   "tasks",
		Short: "List queued tasks",
		Long:  `Usage: tasks`,
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
	root.AddCommand(tasksCmd)

	var aliasCmd = &cobra.Command{
		Use:   "alias",
		Short: "Give an agent alias",
		Long: `Usage: alias <identifier> <alias>
					Valid identifiers are: UUID, Hostname, Username
					Example: alias 192.168.1.1 host1`,
		Run: func(cmd *cobra.Command, args []string) {
			arguments := strings.Split(strings.TrimSpace(strings.TrimPrefix(command, "alias ")), " ")
			if len(arguments) != 2 {
				logger.Info("Invalid identifier, valid identifiers are: UUID, Hostname, Username")
			} else {
				for i, agent := range c2.State.RegisteredAgents {
					if agent.UUID == arguments[0] || agent.Hostname == arguments[0] || agent.Username == arguments[0] {
						c2.State.RegisteredAgents[i].Alias = arguments[1]
						logger.Info("Aliased agent to name: " + c2.State.RegisteredAgents[i].Alias)
						break
					}
					logger.Info("No agents found with given identifier")
				}
			}
		},
	}
	root.AddCommand(aliasCmd)

	var sessionCmd = &cobra.Command{
		Use:   "session",
		Short: "Signal agent to initiate a TCP session",
		Long:  `Usage: session`,
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Initiating session")
			// TODO: Switch between connection types
			go c2.InitiateSession()
		},
	}
	root.AddCommand(sessionCmd)

	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start a teamserver",
		Long:  `Usage: serve`,
		Run: func(cmd *cobra.Command, args []string) {
			if c2.State.IsServing != true {
				c2.State.IsServing = true
				logger.Info("Started teamserver on port: " + config.TeamserverPort)
			} else {
				logger.Info("Teamserver is already running on port: " + config.TeamserverPort)
			}
		},
	}
	root.AddCommand(serveCmd)

	var loglevelCmd = &cobra.Command{
		Use:   "loglevel",
		Short: "Change log level",
		Long: `Usage: loglevel <level>
					Valid log levels are: VERBOSE, INFO, SILENT`,
		Run: func(cmd *cobra.Command, args []string) {
			loglevel := strings.TrimSpace(strings.TrimPrefix(command, "loglevel "))
			switch loglevel {
			case "VERBOSE":
				config.LogLevel = "VERBOSE"
			case "INFO":
				config.LogLevel = "INFO"
			case "SILENT":
				config.LogLevel = "SILENT"
			default:
				logger.Info("Invalid log level, valid log levels are: VERBOSE, INFO, SILENT ")
			}
		},
	}
	root.AddCommand(loglevelCmd)

	var exitCmd = &cobra.Command{
		Use:   "exit",
		Short: "Close and exit server",
		Long:  `Usage: exit`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Graceful exit
			os.Exit(0)
		},
	}
	root.AddCommand(exitCmd)
	return root
}

func logAgent(agent c2.Agent) string {
	return ("UUID: " + agent.UUID + " Alias: " + agent.Alias + " Hostname: " + agent.Hostname)
}
