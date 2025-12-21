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
		strings.TrimSpace(input)
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
			state.Tasks = append(state.Tasks, Task{Command: input})
		}
		mu.Unlock()
	}
}

var commands = []string{"SELECT", "AGENTS", "TASKS"}

func RunCommand(command string) {
	if strings.Contains(command, "SELECT") {
		state.SelectedAgent = strings.TrimPrefix(command, "AGENT ")
		slog.Info("Selected agent with UUID: " + state.SelectedAgent)
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
				slog.Info("Task " + strconv.Itoa(i) + ": " + " Command: " + task.Command)
			}
		}
	} else {
		return
	}
}
