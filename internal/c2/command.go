package c2

import (
	"bufio"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Task struct {
	Agent   string
	Command string
}

var mu sync.Mutex
var task Task

func GetCommands() {
	for {
		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')
		mu.Lock()
		task.Command = strings.TrimSpace(command)
		mu.Unlock()
	}
}

func RunCommands() {
	for {
		mu.Lock()
		if strings.Contains(task.Command, "AGENT") {
			state.selectedAgent = strings.TrimPrefix(task.Command, "AGENT ")
			slog.Info("Selected agent with UUID: " + state.selectedAgent)
			task.Command = ""
		} else if strings.Contains(task.Command, "LIST") {
			if len(state.registeredAgents) == 0 {
				slog.Info("No agents registered")
			} else {
				for i := range len(state.registeredAgents) {
					slog.Info("Agent " + strconv.Itoa(i) + ": " + state.registeredAgents[i])
				}
			}
			task.Command = ""
		}
		mu.Unlock()
	}
}
