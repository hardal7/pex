package main

import (
	"github.com/hardal7/pex/internal/agent"
	"github.com/hardal7/pex/internal/config"
	logger "github.com/hardal7/pex/internal/util"
)

func main() {
	logger.Init()
	config.Load()
	agent.Serve()
}
