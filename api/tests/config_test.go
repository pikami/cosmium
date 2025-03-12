package tests_test

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/pikami/cosmium/api"
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/datastore"
	badgerdatastore "github.com/pikami/cosmium/internal/datastore/badger_datastore"
	mapdatastore "github.com/pikami/cosmium/internal/datastore/map_datastore"
	"github.com/pikami/cosmium/internal/logger"
	"github.com/stretchr/testify/assert"
)

type TestServer struct {
	Server    *httptest.Server
	DataStore datastore.DataStore
	URL       string
}

func getDefaultTestServerConfig() *config.ServerConfig {
	return &config.ServerConfig{
		AccountKey:              config.DefaultAccountKey,
		ExplorerPath:            "/tmp/nothing",
		ExplorerBaseUrlLocation: config.ExplorerBaseUrlLocation,
		DataStore:               "map",
	}
}

func runTestServerCustomConfig(configuration *config.ServerConfig) *TestServer {
	var dataStore datastore.DataStore
	switch configuration.DataStore {
	case config.DataStoreBadger:
		dataStore = badgerdatastore.NewBadgerDataStore()
	default:
		dataStore = mapdatastore.NewMapDataStore(mapdatastore.MapDataStoreOptions{})
	}

	api := api.NewApiServer(dataStore, configuration)

	server := httptest.NewServer(api.GetRouter())

	configuration.DatabaseEndpoint = server.URL

	return &TestServer{
		Server:    server,
		DataStore: dataStore,
		URL:       server.URL,
	}
}

func runTestServer() *TestServer {
	config := getDefaultTestServerConfig()

	config.LogLevel = "debug"
	logger.SetLogLevel(logger.LogLevelDebug)

	return runTestServerCustomConfig(config)
}

const (
	testAccountKey     = "account-key"
	testDatabaseName   = "test-db"
	testCollectionName = "test-coll"
)

type testFunc func(t *testing.T, ts *TestServer, cosmosClient *azcosmos.Client)
type testPreset string

const (
	PresetMapStore    testPreset = "MapDS"
	PresetBadgerStore testPreset = "BadgerDS"
)

func runTestsWithPreset(t *testing.T, name string, testPreset testPreset, f testFunc) {
	serverConfig := getDefaultTestServerConfig()

	serverConfig.LogLevel = "debug"
	logger.SetLogLevel(logger.LogLevelDebug)

	switch testPreset {
	case PresetBadgerStore:
		serverConfig.DataStore = config.DataStoreBadger
	case PresetMapStore:
		serverConfig.DataStore = config.DataStoreMap
	}

	ts := runTestServerCustomConfig(serverConfig)
	defer ts.Server.Close()
	defer ts.DataStore.Close()

	client, err := azcosmos.NewClientFromConnectionString(
		fmt.Sprintf("AccountEndpoint=%s;AccountKey=%s", ts.URL, config.DefaultAccountKey),
		&azcosmos.ClientOptions{},
	)
	assert.Nil(t, err)

	testName := fmt.Sprintf("%s_%s", testPreset, name)

	t.Run(testName, func(t *testing.T) {
		f(t, ts, client)
	})
}

func runTestsWithPresets(t *testing.T, name string, testPresets []testPreset, f testFunc) {
	for _, testPreset := range testPresets {
		runTestsWithPreset(t, name, testPreset, f)
	}
}
