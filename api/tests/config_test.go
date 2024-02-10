package tests_test

import (
	"net/http/httptest"

	"github.com/pikami/cosmium/api"
)

func runTestServer() *httptest.Server {
	return httptest.NewServer(api.CreateRouter())
}

const (
	testAccountKey     = "account-key"
	testDatabaseName   = "test-db"
	testCollectionName = "test-coll"
)
