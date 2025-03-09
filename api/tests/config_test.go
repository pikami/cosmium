package tests_test

import (
	"net/http/httptest"

	"github.com/pikami/cosmium/api"
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/datastore"
	mapdatastore "github.com/pikami/cosmium/internal/datastore/map_datastore"
	"github.com/pikami/cosmium/internal/logger"
)

type TestServer struct {
	Server    *httptest.Server
	DataStore datastore.DataStore
	URL       string
}

func runTestServerCustomConfig(config *config.ServerConfig) *TestServer {
	dataStore := mapdatastore.NewMapDataStore(mapdatastore.MapDataStoreOptions{})

	api := api.NewApiServer(dataStore, config)

	server := httptest.NewServer(api.GetRouter())

	config.DatabaseEndpoint = server.URL

	return &TestServer{
		Server:    server,
		DataStore: dataStore,
		URL:       server.URL,
	}
}

func runTestServer() *TestServer {
	config := &config.ServerConfig{
		AccountKey:              config.DefaultAccountKey,
		ExplorerPath:            "/tmp/nothing",
		ExplorerBaseUrlLocation: config.ExplorerBaseUrlLocation,
	}

	config.LogLevel = "debug"
	logger.SetLogLevel(logger.LogLevelDebug)

	return runTestServerCustomConfig(config)
}

const (
	testAccountKey     = "account-key"
	testDatabaseName   = "test-db"
	testCollectionName = "test-coll"
)
