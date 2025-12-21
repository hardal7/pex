package c2

import (
	"log/slog"
	"net"

	"github.com/hardal7/pex/internal/config"
)

func InitiateSession() {
	listen, err := net.Listen("tcp", config.Host+":"+config.SessionPort)
	if err != nil {
		slog.Info("Failed initiating session")
		return
	} else {
		slog.Info("Initiated session")
	}
	defer listen.Close()
	for {
		connection, _ := listen.Accept()
		go handleRequest(connection)
	}
}

func handleRequest(connection net.Conn) {
	connection.Close()
}
