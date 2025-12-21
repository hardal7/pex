package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const defaultHost string = "localhost"
const defaultPort string = "8080"
const defaultInterval int = 3
const defaultJitter int = 0

var Host = defaultHost
var Port = defaultPort
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
	Port = os.Getenv("PORT")
	Interval, _ = strconv.Atoi(os.Getenv("INTERVAL"))

	slog.Info("Loaded configuration:" + "\nPORT: " + Port + "\nINTERVAL: " + fmt.Sprint(Interval) + "\nJITTER: " + fmt.Sprint(Jitter))
}
