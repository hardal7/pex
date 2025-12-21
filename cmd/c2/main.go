package main

import (
	"github.com/hardal7/pex/internal/c2"
	"github.com/hardal7/pex/internal/config"
)

func main() {
	config.Load()

	go c2.Serve()
	c2.GetCommands()
}
