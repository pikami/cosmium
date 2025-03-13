package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pikami/cosmium/api"
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/datastore"
	badgerdatastore "github.com/pikami/cosmium/internal/datastore/badger_datastore"
	mapdatastore "github.com/pikami/cosmium/internal/datastore/map_datastore"
	"github.com/pikami/cosmium/internal/logger"
)

func main() {
	configuration := config.ParseFlags()

	var dataStore datastore.DataStore
	switch configuration.DataStore {
	case config.DataStoreBadger:
		dataStore = badgerdatastore.NewBadgerDataStore(badgerdatastore.BadgerDataStoreOptions{
			PersistDataFilePath: configuration.PersistDataFilePath,
		})
		logger.InfoLn("Using Badger data store")
	default:
		dataStore = mapdatastore.NewMapDataStore(mapdatastore.MapDataStoreOptions{
			InitialDataFilePath: configuration.InitialDataFilePath,
			PersistDataFilePath: configuration.PersistDataFilePath,
		})
		logger.InfoLn("Using in-memory data store")
	}

	server := api.NewApiServer(dataStore, &configuration)
	err := server.Start()
	if err != nil {
		panic(err)
	}

	waitForExit(server, dataStore)
}

func waitForExit(server *api.ApiServer, dataStore datastore.DataStore) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Block until a exit signal is received
	<-sigs

	// Stop the server
	server.Stop()

	dataStore.Close()
}
