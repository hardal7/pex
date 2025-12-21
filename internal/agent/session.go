package agent

import (
	"log/slog"
	"net"
	"time"

	"github.com/hardal7/pex/internal/config"
)

func JoinSession() {
	time.Sleep(3 * time.Second)
	tcpServer, err := net.ResolveTCPAddr("tcp", config.Host+":"+config.SessionPort)
	connection, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		slog.Info("Failed connecting to session")
		defer connection.Close()
		return
	} else {
		_, err = connection.Write([]byte("Joined session"))
		slog.Info("Join session")
	}
}
