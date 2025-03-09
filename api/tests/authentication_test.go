package tests_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/pikami/cosmium/api/config"
	"github.com/stretchr/testify/assert"
)

func Test_Authentication(t *testing.T) {
	ts := runTestServer()
	defer ts.Server.Close()

	t.Run("Should get 200 when correct account key is used", func(t *testing.T) {
		ts.DataStore.DeleteDatabase(testDatabaseName)
		client, err := azcosmos.NewClientFromConnectionString(
			fmt.Sprintf("AccountEndpoint=%s;AccountKey=%s", ts.URL, config.DefaultAccountKey),
			&azcosmos.ClientOptions{},
		)
		assert.Nil(t, err)

		createResponse, err := client.CreateDatabase(
			context.TODO(),
			azcosmos.DatabaseProperties{ID: testDatabaseName},
			&azcosmos.CreateDatabaseOptions{})
		assert.Nil(t, err)
		assert.Equal(t, createResponse.DatabaseProperties.ID, testDatabaseName)
	})

	t.Run("Should get 401 when wrong account key is used", func(t *testing.T) {
		ts.DataStore.DeleteDatabase(testDatabaseName)
		client, err := azcosmos.NewClientFromConnectionString(
			fmt.Sprintf("AccountEndpoint=%s;AccountKey=%s", ts.URL, "AAAA"),
			&azcosmos.ClientOptions{},
		)
		assert.Nil(t, err)

		_, err = client.CreateDatabase(
			context.TODO(),
			azcosmos.DatabaseProperties{ID: testDatabaseName},
			&azcosmos.CreateDatabaseOptions{})

		assert.Contains(t, err.Error(), "401 Unauthorized")
	})

	t.Run("Should allow unauthorized requests to /_explorer", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/_explorer/config.json")
		assert.Nil(t, err)
		defer res.Body.Close()
		responseBody, err := io.ReadAll(res.Body)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Contains(t, string(responseBody), "BACKEND_ENDPOINT")
	})
}

func Test_Authentication_Disabled(t *testing.T) {
	ts := runTestServerCustomConfig(&config.ServerConfig{
		AccountKey:              config.DefaultAccountKey,
		ExplorerPath:            "/tmp/nothing",
		ExplorerBaseUrlLocation: config.ExplorerBaseUrlLocation,
		DisableAuth:             true,
	})
	defer ts.Server.Close()

	t.Run("Should get 200 when wrong account key is used, but authentication is dissabled", func(t *testing.T) {
		ts.DataStore.DeleteDatabase(testDatabaseName)
		client, err := azcosmos.NewClientFromConnectionString(
			fmt.Sprintf("AccountEndpoint=%s;AccountKey=%s", ts.URL, "AAAA"),
			&azcosmos.ClientOptions{},
		)
		assert.Nil(t, err)

		createResponse, err := client.CreateDatabase(
			context.TODO(),
			azcosmos.DatabaseProperties{ID: testDatabaseName},
			&azcosmos.CreateDatabaseOptions{})
		assert.Nil(t, err)
		assert.Equal(t, createResponse.DatabaseProperties.ID, testDatabaseName)
	})
}
