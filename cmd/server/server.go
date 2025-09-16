package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pikami/cosmium/api"
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/datastore"
	badgerdatastore "github.com/pikami/cosmium/internal/datastore/badger_datastore"
	jsondatastore "github.com/pikami/cosmium/internal/datastore/json_datastore"
	"github.com/pikami/cosmium/internal/logger"
)

func main() {
	configuration := config.ParseFlags()

	var dataStore datastore.DataStore
	switch configuration.DataStore {
	case config.DataStoreBadger:
		dataStore = badgerdatastore.NewBadgerDataStore(badgerdatastore.BadgerDataStoreOptions{
			InitialDataFilePath: configuration.InitialDataFilePath,
			PersistDataFilePath: configuration.PersistDataFilePath,
		})
		logger.InfoLn("Using Badger data store")
	default:
		dataStore = jsondatastore.NewJsonDataStore(jsondatastore.JsonDataStoreOptions{
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
