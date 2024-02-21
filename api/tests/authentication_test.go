package tests_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func Test_Authentication(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()

	t.Run("Should get 200 when correct account key is used", func(t *testing.T) {
		repositories.DeleteDatabase(testDatabaseName)
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

	t.Run("Should get 200 when wrong account key is used, but authentication is dissabled", func(t *testing.T) {
		config.Config.DisableAuth = true
		repositories.DeleteDatabase(testDatabaseName)
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
		config.Config.DisableAuth = false
	})

	t.Run("Should get 401 when wrong account key is used", func(t *testing.T) {
		repositories.DeleteDatabase(testDatabaseName)
		client, err := azcosmos.NewClientFromConnectionString(
			fmt.Sprintf("AccountEndpoint=%s;AccountKey=%s", ts.URL, "AAAA"),
			&azcosmos.ClientOptions{},
		)
		assert.Nil(t, err)

		_, err = client.CreateDatabase(
			context.TODO(),
			azcosmos.DatabaseProperties{ID: testDatabaseName},
			&azcosmos.CreateDatabaseOptions{})

		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) {
			assert.Equal(t, respErr.StatusCode, http.StatusUnauthorized)
		} else {
			panic(err)
		}
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
