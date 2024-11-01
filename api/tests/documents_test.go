package tests_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/repositories"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
	"github.com/stretchr/testify/assert"
)

func testCosmosQuery(t *testing.T,
	collectionClient *azcosmos.ContainerClient,
	query string,
	queryParameters []azcosmos.QueryParameter,
	expectedData []interface{},
) {
	pager := collectionClient.NewQueryItemsPager(
		query,
		azcosmos.PartitionKey{},
		&azcosmos.QueryOptions{
			QueryParameters: queryParameters,
		})

	context := context.TODO()
	items := make([]interface{}, 0)

	for pager.More() {
		response, err := pager.NextPage(context)
		assert.Nil(t, err)

		for _, bytes := range response.Items {
			var item interface{}
			err := json.Unmarshal(bytes, &item)
			assert.Nil(t, err)

			items = append(items, item)
		}
	}

	assert.Equal(t, len(expectedData), len(items))
	if !reflect.DeepEqual(items, expectedData) {
		t.Errorf("executed query does not match expected data.\nExpected: %+v\nGot: %+v", expectedData, items)
	}
}

func documents_InitializeDb(t *testing.T) (*httptest.Server, *azcosmos.ContainerClient) {
	repositories.CreateDatabase(repositorymodels.Database{ID: testDatabaseName})
	repositories.CreateCollection(testDatabaseName, repositorymodels.Collection{
		ID: testCollectionName,
		PartitionKey: struct {
			Paths   []string "json:\"paths\""
			Kind    string   "json:\"kind\""
			Version int      "json:\"Version\""
		}{
			Paths: []string{"/pk"},
		},
	})
	repositories.CreateDocument(testDatabaseName, testCollectionName, map[string]interface{}{"id": "12345", "pk": "123", "isCool": false, "arr": []int{1, 2, 3}})
	repositories.CreateDocument(testDatabaseName, testCollectionName, map[string]interface{}{"id": "67890", "pk": "456", "isCool": true, "arr": []int{6, 7, 8}})

	ts := runTestServer()

	client, err := azcosmos.NewClientFromConnectionString(
		fmt.Sprintf("AccountEndpoint=%s;AccountKey=%s", ts.URL, config.Config.AccountKey),
		&azcosmos.ClientOptions{},
	)
	assert.Nil(t, err)

	collectionClient, err := client.NewContainer(testDatabaseName, testCollectionName)
	assert.Nil(t, err)

	return ts, collectionClient
}

func Test_Documents(t *testing.T) {
	ts, collectionClient := documents_InitializeDb(t)
	defer ts.Close()

	t.Run("Should query document", func(t *testing.T) {
		testCosmosQuery(t, collectionClient,
			"SELECT c.id, c[\"pk\"] FROM c ORDER BY c.id",
			nil,
			[]interface{}{
				map[string]interface{}{"id": "12345", "pk": "123"},
				map[string]interface{}{"id": "67890", "pk": "456"},
			},
		)
	})

	t.Run("Should query VALUE array", func(t *testing.T) {
		testCosmosQuery(t, collectionClient,
			"SELECT VALUE [c.id, c[\"pk\"]] FROM c ORDER BY c.id",
			nil,
			[]interface{}{
				[]interface{}{"12345", "123"},
				[]interface{}{"67890", "456"},
			},
		)
	})

	t.Run("Should query VALUE object", func(t *testing.T) {
		testCosmosQuery(t, collectionClient,
			"SELECT VALUE { id: c.id, _pk: c.pk } FROM c ORDER BY c.id",
			nil,
			[]interface{}{
				map[string]interface{}{"id": "12345", "_pk": "123"},
				map[string]interface{}{"id": "67890", "_pk": "456"},
			},
		)
	})

	t.Run("Should query document with single WHERE condition", func(t *testing.T) {
		testCosmosQuery(t, collectionClient,
			`select c.id
			FROM c
			WHERE c.isCool=true
			ORDER BY c.id`,
			nil,
			[]interface{}{
				map[string]interface{}{"id": "67890"},
			},
		)
	})

	t.Run("Should query document with query parameters", func(t *testing.T) {
		testCosmosQuery(t, collectionClient,
			`select c.id
			FROM c
			WHERE c.id=@param_id
			ORDER BY c.id`,
			[]azcosmos.QueryParameter{
				{Name: "@param_id", Value: "67890"},
			},
			[]interface{}{
				map[string]interface{}{"id": "67890"},
			},
		)
	})

	t.Run("Should query array accessor", func(t *testing.T) {
		testCosmosQuery(t, collectionClient,
			`SELECT c.id,
				c["arr"][0] AS arr0,
				c["arr"][1] AS arr1,
				c["arr"][2] AS arr2,
				c["arr"][3] AS arr3
			FROM c ORDER BY c.id`,
			nil,
			[]interface{}{
				map[string]interface{}{"id": "12345", "arr0": 1.0, "arr1": 2.0, "arr2": 3.0, "arr3": nil},
				map[string]interface{}{"id": "67890", "arr0": 6.0, "arr1": 7.0, "arr2": 8.0, "arr3": nil},
			},
		)
	})
}

func Test_Documents_Patch(t *testing.T) {
	ts, collectionClient := documents_InitializeDb(t)
	defer ts.Close()

	t.Run("Should PATCH document", func(t *testing.T) {
		context := context.TODO()
		expectedData := map[string]interface{}{"id": "67890", "pk": "456", "newField": "newValue"}

		patch := azcosmos.PatchOperations{}
		patch.AppendAdd("/newField", "newValue")
		patch.AppendRemove("/isCool")

		itemResponse, err := collectionClient.PatchItem(
			context,
			azcosmos.PartitionKey{},
			"67890",
			patch,
			&azcosmos.ItemOptions{
				EnableContentResponseOnWrite: false,
			},
		)
		assert.Nil(t, err)

		var itemResponseBody map[string]string
		json.Unmarshal(itemResponse.Value, &itemResponseBody)

		assert.Equal(t, expectedData["id"], itemResponseBody["id"])
		assert.Equal(t, expectedData["pk"], itemResponseBody["pk"])
		assert.Empty(t, itemResponseBody["isCool"])
		assert.Equal(t, expectedData["newField"], itemResponseBody["newField"])
	})

	t.Run("Should not allow to PATCH document ID", func(t *testing.T) {
		context := context.TODO()

		patch := azcosmos.PatchOperations{}
		patch.AppendReplace("/id", "newValue")

		_, err := collectionClient.PatchItem(
			context,
			azcosmos.PartitionKey{},
			"67890",
			patch,
			&azcosmos.ItemOptions{
				EnableContentResponseOnWrite: false,
			},
		)
		assert.NotNil(t, err)

		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) {
			assert.Equal(t, http.StatusUnprocessableEntity, respErr.StatusCode)
		} else {
			panic(err)
		}
	})

	t.Run("CreateItem", func(t *testing.T) {
		context := context.TODO()

		item := map[string]interface{}{
			"Id":       "6789011",
			"pk":       "456",
			"newField": "newValue2",
		}
		bytes, err := json.Marshal(item)
		assert.Nil(t, err)

		r, err2 := collectionClient.CreateItem(
			context,
			azcosmos.PartitionKey{},
			bytes,
			&azcosmos.ItemOptions{
				EnableContentResponseOnWrite: false,
			},
		)
		assert.NotNil(t, r)
		assert.Nil(t, err2)

	})

	t.Run("CreateItem that already exists", func(t *testing.T) {
		context := context.TODO()

		item := map[string]interface{}{"id": "12345", "pk": "123", "isCool": false, "arr": []int{1, 2, 3}}
		bytes, err := json.Marshal(item)
		assert.Nil(t, err)

		r, err2 := collectionClient.CreateItem(
			context,
			azcosmos.PartitionKey{},
			bytes,
			&azcosmos.ItemOptions{
				EnableContentResponseOnWrite: false,
			},
		)
		assert.NotNil(t, r)
		assert.NotNil(t, err2)

	})

	t.Run("UpsertItem new", func(t *testing.T) {
		context := context.TODO()

		item := map[string]interface{}{"id": "123456", "pk": "1234", "isCool": false, "arr": []int{1, 2, 3}}
		bytes, err := json.Marshal(item)
		assert.Nil(t, err)

		r, err2 := collectionClient.UpsertItem(
			context,
			azcosmos.PartitionKey{},
			bytes,
			&azcosmos.ItemOptions{
				EnableContentResponseOnWrite: false,
			},
		)
		assert.NotNil(t, r)
		assert.Nil(t, err2)

	})

	t.Run("UpsertItem that already exists", func(t *testing.T) {
		context := context.TODO()

		item := map[string]interface{}{"id": "12345", "pk": "123", "isCool": false, "arr": []int{1, 2, 3, 4}}
		bytes, err := json.Marshal(item)
		assert.Nil(t, err)

		r, err2 := collectionClient.UpsertItem(
			context,
			azcosmos.PartitionKey{},
			bytes,
			&azcosmos.ItemOptions{
				EnableContentResponseOnWrite: false,
			},
		)
		assert.NotNil(t, r)
		assert.Nil(t, err2)

	})

}
