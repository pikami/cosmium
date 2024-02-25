package main

import (
	"fmt"
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

	router := api.CreateRouter()
	if config.Config.TLS_CertificatePath == "" ||
		config.Config.TLS_CertificateKey == "" {
		go router.Run(fmt.Sprintf(":%d", config.Config.Port))
	} else {
		go router.RunTLS(
			fmt.Sprintf(":%d", config.Config.Port),
			config.Config.TLS_CertificatePath,
			config.Config.TLS_CertificateKey)
	}

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
