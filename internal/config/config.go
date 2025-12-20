package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const defaultPort string = "8080"
const defaultInterval int = 1

var Port = defaultPort
var Interval = defaultInterval

func Load() {
	slog.Info("Loading environment variables")
	err := godotenv.Load()

	if err != nil {
		slog.Error("Failed to load .env variables: " + err.Error())
		return
	}

	Port = os.Getenv("PORT")
	Interval, _ = strconv.Atoi(os.Getenv("INTERVAL"))

	slog.Info("Loaded configuration:" + "\nPORT: " + Port + "\nINTERVAL: " + fmt.Sprint(Interval))
}
