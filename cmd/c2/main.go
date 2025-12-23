package main

import (
	"github.com/hardal7/pex/internal/c2"
	"github.com/hardal7/pex/internal/config"
	"github.com/hardal7/pex/internal/console"
	logger "github.com/hardal7/pex/internal/util"
)

func main() {
	logger.Load()
	config.Load()
	go console.RunApp()
	c2.Run()
}
