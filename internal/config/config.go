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

var Host = defaultHost
var BeaconPort = defaultBeaconPort
var SessionPort = defaultSessionPort
var Interval = defaultInterval
var Jitter = defaultJitter

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

	slog.Info("Loaded configuration:" + "\nBEACON PORT: " + BeaconPort + "\nSESSION PORT: " + SessionPort + "\nINTERVAL: " + fmt.Sprint(Interval) + "\nJITTER: " + fmt.Sprint(Jitter))
}
