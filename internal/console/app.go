package console

import (
	"fmt"

	logger "github.com/hardal7/pex/internal/util"
	"github.com/reeflective/console"
)

var ConsoleApp *console.Console

func RunApp() {

	ConsoleApp = console.New("[pex]")

	ConsoleApp.SetPrintLogo(func(_ *console.Console) {
		fmt.Print(`
`)
	})

	menu := ConsoleApp.ActiveMenu()
	menu.SetCommands(MenuCommands)
	InitCommands()
	err := ConsoleApp.Start()
	if err != nil {
		logger.Error("Failed to create console application")
	}
}
