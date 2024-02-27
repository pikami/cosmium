package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pikami/cosmium/api"
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/repositories"
)

func main() {
	config.ParseFlags()

	if config.Config.InitialDataFilePath != "" {
		repositories.LoadStateFS(config.Config.InitialDataFilePath)
	}

	go api.StartAPI()

	waitForExit()
}

func waitForExit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Block until a exit signal is received
	<-sigs

	if config.Config.PersistDataFilePath != "" {
		repositories.SaveStateFS(config.Config.PersistDataFilePath)
	}
}
