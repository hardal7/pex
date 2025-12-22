package c2

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"sync"

	logger "github.com/hardal7/pex/internal/util"
)

var mu sync.Mutex

func GetCommands() {
	for {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		mu.Lock()
		isShellCommand := true
		for _, command := range commands {
			if strings.Contains(input, command) {
				go RunCommand(input)
				isShellCommand = false
				break
			}
		}
		if isShellCommand {
			if state.SelectedAgent.UUID == "ALL" {
				for _, agent := range state.RegisteredAgents {
					state.Tasks = append(state.Tasks, Task{Command: input, Recipient: agent})
				}
			} else if state.SelectedAgent.UUID != "NONE" {
				state.Tasks = append(state.Tasks, Task{Command: input, Recipient: Agent{UUID: state.SelectedAgent.UUID}})
			}
		}
		mu.Unlock()
	}
}

var commands = []string{"SELECT", "AGENTS", "TASKS", "ALIAS"}

const uuidLength int = 32

func RunCommand(command string) {
	if strings.Contains(command, "SELECT") {
		mu.Lock()

		selectName := strings.TrimSpace(strings.TrimPrefix(command, "SELECT "))
		if len(selectName) == uuidLength {
			state.SelectedAgent.UUID = selectName
		} else if strings.Contains(selectName, ":") || strings.Contains(selectName, ".") {
			state.SelectedAgent.Hostname = selectName
		} else {
			state.SelectedAgent.Alias = selectName
		}
		if selectName == "ALL" {
			logger.Info("Selected all agents")
		} else {
			for _, agent := range state.RegisteredAgents {
				if state.SelectedAgent.UUID == agent.UUID || state.SelectedAgent.Hostname == agent.Hostname || state.SelectedAgent.Alias == agent.Alias {
					logger.Info("Selected agent with " + logAgent(state.SelectedAgent))
					break
				} else {
					logger.Info("No registered agent found with UUID: " + state.SelectedAgent.UUID)
					state.SelectedAgent.UUID = "NONE"
				}
			}
		}
		mu.Unlock()
	} else if strings.Contains(command, "AGENTS") {
		if len(state.RegisteredAgents) == 0 {
			logger.Info("No agents registered")
		} else {
			for i := range len(state.RegisteredAgents) {
				logger.Info("Agent " + strconv.Itoa(i) + ": " + logAgent(state.RegisteredAgents[i]))
			}
		}
	} else if strings.Contains(command, "TASKS") {
		if len(state.Tasks) == 0 {
			logger.Info("No tasks in queue")
		} else {
			for i, task := range state.Tasks {
				logger.Info("Task " + strconv.Itoa(i) + ": " + " Command: " + task.Command + " Recipient: " + task.Recipient.UUID)
			}
		}
	} else if strings.Contains(command, "ALIAS") {
		arguments := strings.Split(strings.TrimSpace(strings.TrimPrefix(command, "ALIAS ")), " ")
		if len(arguments) != 2 {
			logger.Info("Invalid identifier, valid identifiers are: UUID, Hostname, Username")
			logger.Info("Command ALIAS Usage: ALIAS <identifier> <alias>")
		} else {
			for i, agent := range state.RegisteredAgents {
				if agent.UUID == arguments[0] || agent.Hostname == arguments[0] || agent.Username == arguments[0] {
					state.RegisteredAgents[i].Alias = arguments[1]
					logger.Info("Aliased agent to name: " + state.RegisteredAgents[i].Alias)
					break
				}
				logger.Info("No agents found with given identifier")
			}
		}
	} else {
		return
	}
}

func logAgent(agent Agent) string {
	return ("UUID: " + agent.UUID + " Alias: " + agent.Alias + " Hostname: " + agent.Hostname)
}
