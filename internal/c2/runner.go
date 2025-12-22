package c2

import (
	"log/slog"

	"github.com/hardal7/pex/internal/config"
)

type Task struct {
	Agent   string
	Command string
}

type ServerState struct {
	RegisteredAgents []string
	SelectedAgent    string
	Tasks            []Task
}

var state ServerState = ServerState{SelectedAgent: "NONE"}

func Run() {
	go GetCommands()

	switch config.ConnectionType {
	case "http":
		ServeHTTP()
	case "tcp":
		InitiateSession()
	default:
		slog.Error("Invalid connection type: " + config.ConnectionType)
	}
}
