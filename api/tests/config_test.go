package tests_test

import (
	"net/http/httptest"

	"github.com/pikami/cosmium/api"
	"github.com/pikami/cosmium/api/config"
)

func runTestServer() *httptest.Server {
	config.Config.AccountKey = config.DefaultAccountKey
	config.Config.ExplorerPath = "/tmp/nothing"

	return httptest.NewServer(api.CreateRouter())
}

const (
	testAccountKey     = "account-key"
	testDatabaseName   = "test-db"
	testCollectionName = "test-coll"
)
