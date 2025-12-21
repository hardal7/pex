package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const defaultHost string = "localhost"
const defaultBeaconPort string = "8080"
const defaultSessionPort string = "9090"
const defaultInterval int = 3
const defaultJitter int = 0
const defaultConnectionType string = "HTTP"

var Host = defaultHost
var BeaconPort = defaultBeaconPort
var SessionPort = defaultSessionPort
var Interval = defaultInterval
var Jitter = defaultJitter
var ConnectionType = defaultConnectionType

func Load() {
	slog.Info("Loading environment variables")
	err := godotenv.Load()

	if err != nil {
		slog.Error("Failed to load .env variables: " + err.Error())
		return
	}

	Host = os.Getenv("HOST")
	BeaconPort = os.Getenv("BEACON_PORT")
	SessionPort = os.Getenv("SESSION_PORT")
	Interval, _ = strconv.Atoi(os.Getenv("INTERVAL"))
	Jitter, _ = strconv.Atoi(os.Getenv("JITTER"))
	ConnectionType = os.Getenv("CONNECTION_TYPE")

	slog.Info("Loaded configuration:" + "\nBEACON PORT: " + BeaconPort + "\nSESSION PORT: " + SessionPort + "\nINTERVAL: " + fmt.Sprint(Interval) + "\nJITTER: " + fmt.Sprint(Jitter) + "\n CONNECTION TYPE: " + ConnectionType)
}
