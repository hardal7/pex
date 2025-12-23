package console

import (
	"github.com/fatih/color"
	logger "github.com/hardal7/pex/internal/util"
	"github.com/reeflective/console"
)

var ConsoleApp *console.Console

func RunApp() {

	ConsoleApp = console.New("[pex]")

	ConsoleApp.SetPrintLogo(func(_ *console.Console) {
		color.Red(`
	██████╗ ███████╗██╗  ██╗
	██╔══██╗██╔════╝╚██╗██╔╝
	██████╔╝█████╗   ╚███╔╝ 
	██╔═══╝ ██╔══╝   ██╔██╗ 
	██║     ███████╗██╔╝ ██╗
	╚═╝     ╚══════╝╚═╝  ╚═╝
						
	Made by hardal with <3
						
	> github.com/hardal7
	> https://pex.sh/docs

	(C) ISC LICENSE 2025

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
