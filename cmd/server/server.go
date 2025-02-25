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
	configuration := config.ParseFlags()

	repository := repositories.NewDataRepository(repositories.RepositoryOptions{
		InitialDataFilePath: configuration.InitialDataFilePath,
		PersistDataFilePath: configuration.PersistDataFilePath,
	})

	server := api.NewApiServer(repository, &configuration)
	err := server.Start()
	if err != nil {
		panic(err)
	}

	waitForExit(server, repository, configuration)
}

func waitForExit(server *api.ApiServer, repository *repositories.DataRepository, config config.ServerConfig) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Block until a exit signal is received
	<-sigs

	// Stop the server
	server.Stop()

	if config.PersistDataFilePath != "" {
		repository.SaveStateFS(config.PersistDataFilePath)
	}
}
