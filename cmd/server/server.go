package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pikami/cosmium/api"
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/datastore"
	mapdatastore "github.com/pikami/cosmium/internal/datastore/map_datastore"
)

func main() {
	configuration := config.ParseFlags()

	var dataStore datastore.DataStore = mapdatastore.NewMapDataStore(mapdatastore.MapDataStoreOptions{
		InitialDataFilePath: configuration.InitialDataFilePath,
		PersistDataFilePath: configuration.PersistDataFilePath,
	})

	server := api.NewApiServer(dataStore, &configuration)
	err := server.Start()
	if err != nil {
		panic(err)
	}

	waitForExit(server, dataStore, configuration)
}

func waitForExit(server *api.ApiServer, dataStore datastore.DataStore, config config.ServerConfig) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Block until a exit signal is received
	<-sigs

	// Stop the server
	server.Stop()

	dataStore.Close()
}
