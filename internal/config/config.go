package config

import (
	"fmt"
	"log/slog"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const defaultHost string = "localhost"
const defaultBeaconPort string = "8080"
const defaultSessionPort string = "9090"
const defaultInterval int = 3
const defaultJitter int = 0
const defaultConnectionType string = "http"

var Host = defaultHost
var BeaconPort = defaultBeaconPort
var SessionPort = defaultSessionPort
var Interval = defaultInterval
var Jitter = defaultJitter
var ConnectionType = defaultConnectionType

func Load() {
	slog.Info("Loading environment variables")
	var k = koanf.New(".")

	if err := k.Load(file.Provider("config/config.yml"), yaml.Parser()); err != nil {
		slog.Info("No config file found, using defaults")
	} else {
		Host = k.String("host")
		BeaconPort = k.String("port.beacon")
		SessionPort = k.String("port.session")
		Interval = k.Int("beacon.interval")
		Jitter = k.Int("beacon.jitter")
		ConnectionType = k.String("connection")
	}

	slog.Info("Loaded configuration:" + "\nBEACON PORT: " + BeaconPort + "\nSESSION PORT: " + SessionPort + "\nINTERVAL: " + fmt.Sprint(Interval) + "\nJITTER: " + fmt.Sprint(Jitter) + "\nCONNECTION TYPE: " + ConnectionType)
}
