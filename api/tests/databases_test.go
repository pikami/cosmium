package tests_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/stretchr/testify/assert"
)

func Test_Databases(t *testing.T) {
	presets := []testPreset{PresetJsonStore, PresetBadgerStore}

	runTestsWithPresets(t, "Database Create", presets, func(t *testing.T, ts *TestServer, client *azcosmos.Client) {
		t.Run("Should create database", func(t *testing.T) {
			ts.DataStore.DeleteDatabase(testDatabaseName)

			createResponse, err := client.CreateDatabase(context.TODO(), azcosmos.DatabaseProperties{
				ID: testDatabaseName,
			}, &azcosmos.CreateDatabaseOptions{})

			assert.Nil(t, err)
			assert.Equal(t, createResponse.DatabaseProperties.ID, testDatabaseName)
		})

		t.Run("Should return conflict when database exists", func(t *testing.T) {
			ts.DataStore.CreateDatabase(datastore.Database{
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

	runTestsWithPresets(t, "Database Read", presets, func(t *testing.T, ts *TestServer, client *azcosmos.Client) {
		t.Run("Should read database", func(t *testing.T) {
			ts.DataStore.CreateDatabase(datastore.Database{
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
			ts.DataStore.DeleteDatabase(testDatabaseName)

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

	runTestsWithPresets(t, "Database Delete", presets, func(t *testing.T, ts *TestServer, client *azcosmos.Client) {
		t.Run("Should delete database", func(t *testing.T) {
			ts.DataStore.CreateDatabase(datastore.Database{
				ID: testDatabaseName,
			})

			databaseResponse, err := client.NewDatabase(testDatabaseName)
			assert.Nil(t, err)

			readResponse, err := databaseResponse.Delete(context.TODO(), &azcosmos.DeleteDatabaseOptions{})
			assert.Nil(t, err)
			assert.Equal(t, readResponse.RawResponse.StatusCode, http.StatusNoContent)
		})

		t.Run("Should return not found when database does not exist", func(t *testing.T) {
			ts.DataStore.DeleteDatabase(testDatabaseName)

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

		t.Run("Should delete database with exactly matching name", func(t *testing.T) {
			ts.DataStore.CreateDatabase(datastore.Database{
				ID: testDatabaseName + "extra",
			})
			ts.DataStore.CreateDatabase(datastore.Database{
				ID: testDatabaseName,
			})

			databaseResponse, err := client.NewDatabase(testDatabaseName)
			assert.Nil(t, err)

			readResponse, err := databaseResponse.Delete(context.TODO(), &azcosmos.DeleteDatabaseOptions{})
			assert.Nil(t, err)
			assert.Equal(t, readResponse.RawResponse.StatusCode, http.StatusNoContent)

			dbs, status := ts.DataStore.GetAllDatabases()
			assert.Equal(t, status, datastore.StatusOk)
			assert.Len(t, dbs, 1)
			assert.Equal(t, dbs[0].ID, testDatabaseName+"extra")
		})
	})
}
