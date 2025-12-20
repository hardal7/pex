package main

import (
	"github.com/hardal7/pex/internal/c2"
	"github.com/hardal7/pex/internal/config"
)

func main() {
	config.Load()

	ch := make(chan string)
	go c2.Serve(ch)
	c2.GetCommands(ch)
}
