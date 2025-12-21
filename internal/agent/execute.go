package agent

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

func ExecuteCommand(command []string) string {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	var commandExists bool = false
	var cmd *exec.Cmd
	if len(command) != 0 {
		command = append([]string{"-c"}, strings.Join(command, " "))
		cmd = exec.Command("bash", command...)
		commandExists = true
	}

	if commandExists {
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			slog.Error("Failed executing command: " + err.Error())
			fmt.Println(command)
		}
	}
	if stderr.String() != "" {
		slog.Info("Error Output: " + stderr.String())
	}

	return stdout.String() + stderr.String()
}
