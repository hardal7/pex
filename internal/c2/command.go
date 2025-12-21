package c2

import (
	"bufio"
	"log/slog"
	"os"
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
		command, err := reader.ReadString('\n')

		if err != nil {
			slog.Info("Error reading command: " + err.Error())
		} else {
			mu.Lock()
			task.Command = command
			mu.Unlock()
		}
	}
}
