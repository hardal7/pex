package console

import (
	"log/slog"

	"github.com/fatih/color"
	"github.com/hardal7/pex/internal/c2"
	"github.com/reeflective/console"
)

var ConsoleApp *console.Console

func RunApp() {

	// TODO: Add selected agent prefix to the prompt
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
	menu.SetCommands(c2.MenuCommands)
	c2.InitCommands()
	err := ConsoleApp.Start()
	if err != nil {
		slog.Error("Failed to create logger application")
	}
}
