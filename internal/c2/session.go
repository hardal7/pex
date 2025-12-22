package c2

import (
	"net"

	"github.com/hardal7/pex/internal/config"
	logger "github.com/hardal7/pex/internal/util"
)

func InitiateSession() {
	listen, err := net.Listen("tcp", config.Host+":"+config.SessionPort)
	if err != nil {
		logger.Info("Failed initiating session")
		return
	} else {
		logger.Info("Initiated session")
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
