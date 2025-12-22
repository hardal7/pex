package c2

import (
	"github.com/hardal7/pex/internal/config"
	logger "github.com/hardal7/pex/internal/util"
)

var State ServerState = ServerState{SelectedAgent: Agent{Alias: "NONE"}}

func Run() {
	// go GetCommands()

	switch config.ConnectionType {
	case "http":
		ServeHTTP()
	case "tcp":
		InitiateSession()
	default:
		logger.Error("Invalid connection type: " + config.ConnectionType)
	}
}
