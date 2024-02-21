package tests_test

import (
	"net/http/httptest"

	"github.com/pikami/cosmium/api"
	"github.com/pikami/cosmium/api/config"
)

func runTestServer() *httptest.Server {
	config.Config.AccountKey = config.DefaultAccountKey

	return httptest.NewServer(api.CreateRouter())
}

const (
	testAccountKey     = "account-key"
	testDatabaseName   = "test-db"
	testCollectionName = "test-coll"
)
