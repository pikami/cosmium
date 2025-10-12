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

func Test_Collections(t *testing.T) {
	presets := []testPreset{PresetJsonStore, PresetBadgerStore}

	setUp := func(ts *TestServer, client *azcosmos.Client) *azcosmos.DatabaseClient {
		ts.DataStore.CreateDatabase(datastore.Database{ID: testDatabaseName})
		databaseClient, err := client.NewDatabase(testDatabaseName)
		assert.Nil(t, err)

		return databaseClient
	}

	runTestsWithPresets(t, "Collection Create", presets, func(t *testing.T, ts *TestServer, client *azcosmos.Client) {
		databaseClient := setUp(ts, client)

		t.Run("Should create collection", func(t *testing.T) {
			createResponse, err := databaseClient.CreateContainer(context.TODO(), azcosmos.ContainerProperties{
				ID: testCollectionName,
			}, &azcosmos.CreateContainerOptions{})

			assert.Nil(t, err)
			assert.Equal(t, createResponse.ContainerProperties.ID, testCollectionName)
		})

		t.Run("Should return conflict when collection exists", func(t *testing.T) {
			ts.DataStore.CreateCollection(testDatabaseName, datastore.Collection{
				ID: testCollectionName,
			})

			_, err := databaseClient.CreateContainer(context.TODO(), azcosmos.ContainerProperties{
				ID: testCollectionName,
			}, &azcosmos.CreateContainerOptions{})
			assert.NotNil(t, err)

			var respErr *azcore.ResponseError
			if errors.As(err, &respErr) {
				assert.Equal(t, respErr.StatusCode, http.StatusConflict)
			} else {
				panic(err)
			}
		})
	})

	runTestsWithPresets(t, "Collection Read", presets, func(t *testing.T, ts *TestServer, client *azcosmos.Client) {
		databaseClient := setUp(ts, client)

		t.Run("Should read collection", func(t *testing.T) {
			ts.DataStore.CreateCollection(testDatabaseName, datastore.Collection{
				ID: testCollectionName,
			})

			collectionResponse, err := databaseClient.NewContainer(testCollectionName)
			assert.Nil(t, err)

			readResponse, err := collectionResponse.Read(context.TODO(), &azcosmos.ReadContainerOptions{})
			assert.Nil(t, err)
			assert.Equal(t, readResponse.RawResponse.StatusCode, http.StatusOK)
			assert.Equal(t, readResponse.ContainerProperties.ID, testCollectionName)
		})

		t.Run("Should return not found when collection does not exist", func(t *testing.T) {
			ts.DataStore.DeleteCollection(testDatabaseName, testCollectionName)

			collectionResponse, err := databaseClient.NewContainer(testCollectionName)
			assert.Nil(t, err)

			_, err = collectionResponse.Read(context.TODO(), &azcosmos.ReadContainerOptions{})
			assert.NotNil(t, err)

			var respErr *azcore.ResponseError
			if errors.As(err, &respErr) {
				assert.Equal(t, respErr.StatusCode, http.StatusNotFound)
			} else {
				panic(err)
			}
		})
	})

	runTestsWithPresets(t, "Collection Delete", presets, func(t *testing.T, ts *TestServer, client *azcosmos.Client) {
		databaseClient := setUp(ts, client)

		t.Run("Should delete collection", func(t *testing.T) {
			ts.DataStore.CreateCollection(testDatabaseName, datastore.Collection{
				ID: testCollectionName,
			})

			collectionResponse, err := databaseClient.NewContainer(testCollectionName)
			assert.Nil(t, err)

			readResponse, err := collectionResponse.Delete(context.TODO(), &azcosmos.DeleteContainerOptions{})
			assert.Nil(t, err)
			assert.Equal(t, readResponse.RawResponse.StatusCode, http.StatusNoContent)
		})

		t.Run("Should return not found when collection does not exist", func(t *testing.T) {
			ts.DataStore.DeleteCollection(testDatabaseName, testCollectionName)

			collectionResponse, err := databaseClient.NewContainer(testCollectionName)
			assert.Nil(t, err)

			_, err = collectionResponse.Read(context.TODO(), &azcosmos.ReadContainerOptions{})
			assert.NotNil(t, err)

			var respErr *azcore.ResponseError
			if errors.As(err, &respErr) {
				assert.Equal(t, respErr.StatusCode, http.StatusNotFound)
			} else {
				panic(err)
			}
		})

		t.Run("Should delete collection with exactly matching name", func(t *testing.T) {
			ts.DataStore.CreateCollection(testDatabaseName, datastore.Collection{
				ID: testCollectionName + "extra",
			})
			ts.DataStore.CreateCollection(testDatabaseName, datastore.Collection{
				ID: testCollectionName,
			})

			collectionResponse, err := databaseClient.NewContainer(testCollectionName)
			assert.Nil(t, err)

			readResponse, err := collectionResponse.Delete(context.TODO(), &azcosmos.DeleteContainerOptions{})
			assert.Nil(t, err)
			assert.Equal(t, readResponse.RawResponse.StatusCode, http.StatusNoContent)

			collections, status := ts.DataStore.GetAllCollections(testDatabaseName)
			assert.Equal(t, status, datastore.StatusOk)
			assert.Len(t, collections, 1)
			assert.Equal(t, collections[0].ID, testCollectionName+"extra")
		})
	})
}
