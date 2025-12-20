package c2

import (
	"bufio"
	"log/slog"
	"os"
)

func GetCommands(ch chan string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		command, err := reader.ReadString('\n')

		if err != nil {
			slog.Info("Error reading command: " + err.Error())
		} else {
			ch <- command
		}
	}
}
