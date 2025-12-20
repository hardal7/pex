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
const defaultInterval int = 1
const defaultJitter int = 1
const defaultKeylogBlockSize int = 32

var Host = defaultHost
var Port = defaultPort
var Interval = defaultInterval
var Jitter = defaultJitter
var KeyLogBlockSize = defaultKeylogBlockSize

func Load() {
	slog.Info("Loading environment variables")
	err := godotenv.Load()

	if err != nil {
		slog.Error("Failed to load .env variables: " + err.Error())
		return
	}

	Port = os.Getenv("HOST")
	Port = os.Getenv("PORT")
	Interval, _ = strconv.Atoi(os.Getenv("INTERVAL"))
	KeyLogBlockSize, _ = strconv.Atoi(os.Getenv("KEYLOG_BLOCK_SIZE"))

	slog.Info("Loaded configuration:" + "\nPORT: " + Port + "\nINTERVAL: " + fmt.Sprint(Interval) + "\nJITTER: " + fmt.Sprint(Jitter) + "\nKEYLOG BLOCK SIZE: " + fmt.Sprint(KeyLogBlockSize))
}
