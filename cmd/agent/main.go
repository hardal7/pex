package main

import (
	"github.com/hardal7/pex/internal/agent"
	"github.com/hardal7/pex/internal/config"
)

func main() {
	config.Load()
	agent.Serve()
}
