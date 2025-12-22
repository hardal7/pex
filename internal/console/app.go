package console

import (
	"fmt"
	"log/slog"

	"github.com/reeflective/console"
)

var ConsoleApp *console.Console
var IsRunning bool

func RunApp() {

	ConsoleApp = console.New("[pex]")

	ConsoleApp.SetPrintLogo(func(_ *console.Console) {
		fmt.Print(`
`)
	})

	IsRunning = true
	menu := ConsoleApp.ActiveMenu()
	menu.SetCommands(Commands)
	err := ConsoleApp.Start()
	if err != nil {
		slog.Error("Failed to create console application")
	}
}
