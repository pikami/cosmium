package tests_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/repositories"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
	"github.com/stretchr/testify/assert"
)

func Test_Collections(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()

	client, err := azcosmos.NewClientFromConnectionString(
		fmt.Sprintf("AccountEndpoint=%s;AccountKey=%s", ts.URL, config.Config.AccountKey),
		&azcosmos.ClientOptions{},
	)
	assert.Nil(t, err)

	databaseClient, err := client.NewDatabase(testDatabaseName)
	assert.Nil(t, err)

	t.Run("Collection Create", func(t *testing.T) {
		t.Run("Should create collection", func(t *testing.T) {
			createResponse, err := databaseClient.CreateContainer(context.TODO(), azcosmos.ContainerProperties{
				ID: testCollectionName,
			}, &azcosmos.CreateContainerOptions{})

			assert.Nil(t, err)
			assert.Equal(t, createResponse.ContainerProperties.ID, testCollectionName)
		})

		t.Run("Should return conflict when collection exists", func(t *testing.T) {
			repositories.CreateCollection(testDatabaseName, repositorymodels.Collection{
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

	t.Run("Collection Read", func(t *testing.T) {
		t.Run("Should read collection", func(t *testing.T) {
			repositories.CreateCollection(testDatabaseName, repositorymodels.Collection{
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
			repositories.DeleteCollection(testDatabaseName, testCollectionName)

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

	t.Run("Collection Delete", func(t *testing.T) {
		t.Run("Should delete collection", func(t *testing.T) {
			repositories.CreateCollection(testDatabaseName, repositorymodels.Collection{
				ID: testCollectionName,
			})

			collectionResponse, err := databaseClient.NewContainer(testCollectionName)
			assert.Nil(t, err)

			readResponse, err := collectionResponse.Delete(context.TODO(), &azcosmos.DeleteContainerOptions{})
			assert.Nil(t, err)
			assert.Equal(t, readResponse.RawResponse.StatusCode, http.StatusNoContent)
		})

		t.Run("Should return not found when collection does not exist", func(t *testing.T) {
			repositories.DeleteCollection(testDatabaseName, testCollectionName)

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
}
