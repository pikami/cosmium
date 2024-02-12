package tests_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/pikami/cosmium/internal/repositories"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
	"github.com/stretchr/testify/assert"
)

func Test_Databases(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()

	client, err := azcosmos.NewClientFromConnectionString(
		fmt.Sprintf("AccountEndpoint=%s;AccountKey=%s", ts.URL, "asas"),
		&azcosmos.ClientOptions{},
	)
	assert.Nil(t, err)

	t.Run("Database Create", func(t *testing.T) {
		t.Run("Should create database", func(t *testing.T) {
			createResponse, err := client.CreateDatabase(context.TODO(), azcosmos.DatabaseProperties{
				ID: testDatabaseName,
			}, &azcosmos.CreateDatabaseOptions{})

			assert.Nil(t, err)
			assert.Equal(t, createResponse.DatabaseProperties.ID, testDatabaseName)
		})

		t.Run("Should return conflict when database exists", func(t *testing.T) {
			repositories.CreateDatabase(repositorymodels.Database{
				ID: testDatabaseName,
			})

			_, err := client.CreateDatabase(context.TODO(), azcosmos.DatabaseProperties{
				ID: testDatabaseName,
			}, &azcosmos.CreateDatabaseOptions{})
			assert.NotNil(t, err)

			var respErr *azcore.ResponseError
			if errors.As(err, &respErr) {
				assert.Equal(t, respErr.StatusCode, http.StatusConflict)
			} else {
				panic(err)
			}
		})
	})

	t.Run("Database Read", func(t *testing.T) {
		t.Run("Should read database", func(t *testing.T) {
			repositories.CreateDatabase(repositorymodels.Database{
				ID: testDatabaseName,
			})

			databaseResponse, err := client.NewDatabase(testDatabaseName)
			assert.Nil(t, err)

			readResponse, err := databaseResponse.Read(context.TODO(), &azcosmos.ReadDatabaseOptions{})
			assert.Nil(t, err)
			assert.Equal(t, readResponse.RawResponse.StatusCode, http.StatusOK)
			assert.Equal(t, readResponse.DatabaseProperties.ID, testDatabaseName)
		})

		t.Run("Should return not found when database does not exist", func(t *testing.T) {
			repositories.DeleteDatabase(testDatabaseName)

			databaseResponse, err := client.NewDatabase(testDatabaseName)
			assert.Nil(t, err)

			_, err = databaseResponse.Read(context.TODO(), &azcosmos.ReadDatabaseOptions{})
			assert.NotNil(t, err)

			var respErr *azcore.ResponseError
			if errors.As(err, &respErr) {
				assert.Equal(t, respErr.StatusCode, http.StatusNotFound)
			} else {
				panic(err)
			}
		})
	})

	t.Run("Database Delete", func(t *testing.T) {
		t.Run("Should delete database", func(t *testing.T) {
			repositories.CreateDatabase(repositorymodels.Database{
				ID: testDatabaseName,
			})

			databaseResponse, err := client.NewDatabase(testDatabaseName)
			assert.Nil(t, err)

			readResponse, err := databaseResponse.Delete(context.TODO(), &azcosmos.DeleteDatabaseOptions{})
			assert.Nil(t, err)
			assert.Equal(t, readResponse.RawResponse.StatusCode, http.StatusNoContent)
		})

		t.Run("Should return not found when database does not exist", func(t *testing.T) {
			repositories.DeleteDatabase(testDatabaseName)

			databaseResponse, err := client.NewDatabase(testDatabaseName)
			assert.Nil(t, err)

			_, err = databaseResponse.Read(context.TODO(), &azcosmos.ReadDatabaseOptions{})
			assert.NotNil(t, err)

			var respErr *azcore.ResponseError
			if errors.As(err, &respErr) {
				assert.Equal(t, respErr.StatusCode, http.StatusNotFound)
			} else {
				panic(err)
			}
		})
	})
}
