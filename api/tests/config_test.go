package tests_test

import (
	"net/http/httptest"

	"github.com/pikami/cosmium/api"
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/repositories"
)

type TestServer struct {
	Server     *httptest.Server
	Repository *repositories.DataRepository
	URL        string
}

func runTestServerCustomConfig(config config.ServerConfig) *TestServer {
	repository := repositories.NewDataRepository(repositories.RepositoryOptions{})

	api := api.NewApiServer(repository, config)

	server := httptest.NewServer(api.GetRouter())

	return &TestServer{
		Server:     server,
		Repository: repository,
		URL:        server.URL,
	}
}

func runTestServer() *TestServer {
	config := config.ServerConfig{
		AccountKey:              config.DefaultAccountKey,
		ExplorerPath:            "/tmp/nothing",
		ExplorerBaseUrlLocation: config.ExplorerBaseUrlLocation,
	}

	return runTestServerCustomConfig(config)
}

const (
	testAccountKey     = "account-key"
	testDatabaseName   = "test-db"
	testCollectionName = "test-coll"
)
