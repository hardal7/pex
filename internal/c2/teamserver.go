package c2

import (
	"net"
	"strings"

	"github.com/hardal7/pex/internal/config"
	logger "github.com/hardal7/pex/internal/util"
)

func HostTeamserver() {
	State.IsServing = true
	listen, err := net.Listen("tcp", config.Host+":"+config.TeamserverPort)
	if err != nil {
		logger.Info("Failed starting teamserver: " + err.Error())
		State.IsServing = false
		return
	}
	logger.Info("Started teamserver on port: " + config.TeamserverPort)
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			logger.Info("Invalid connection:" + err.Error())
			return
		}
		go teamserverRequestHandler(conn)
	}
}

func teamserverRequestHandler(conn net.Conn) {
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	defer conn.Close()
	if err != nil {
		logger.Info("Failed reading message from client: " + err.Error())
		return
	} else {
		logger.Info("Received message from client: " + string(buffer))
	}

	command := strings.Split(strings.TrimRight(string(buffer), " "), "\x00")
	if len(command) > 1 {
		ExecuteCommand(command[0], command[2:])
	} else {
		ExecuteCommand(command[0], nil)
	}
}
