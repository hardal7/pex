package config

import (
	"fmt"
	"log/slog"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var Host string = "192.168.1.37"
var BeaconPort string = "8080"
var SessionPort string = "9090"
var TeamserverPort string = "7070"
var Interval int = 3
var Jitter int = 0
var ConnectionType string = "http"
var LogLevel string = "debug"

func Load() {
	slog.Info("Loading environment variables")
	var k = koanf.New(".")

	// FIXME: File provider doesn't work for windows
	if err := k.Load(file.Provider("config/config.yml"), yaml.Parser()); err != nil {
		slog.Info("No config file found, using defaults")
	} else {
		Host = k.String("host")
		BeaconPort = k.String("port.beacon")
		SessionPort = k.String("port.session")
		TeamserverPort = k.String("port.server")
		Interval = k.Int("beacon.interval")
		Jitter = k.Int("beacon.jitter")
		ConnectionType = k.String("connection")
		LogLevel = k.String("loglevel")
	}

	slog.Info("Loaded configuration:" + "\nHOST: " + Host + "\nBEACON PORT: " + BeaconPort + "\nSESSION PORT: " + SessionPort + "\nTEAM SERVER PORT: " + TeamserverPort + "\nINTERVAL: " + fmt.Sprint(Interval) + "\nJITTER: " + fmt.Sprint(Jitter) + "\nCONNECTION TYPE: " + ConnectionType + "\nLOG LEVEL: " + LogLevel)
}
