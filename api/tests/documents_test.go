package tests_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/pikami/cosmium/api/config"
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

func documents_InitializeDb(t *testing.T) (*TestServer, *azcosmos.ContainerClient) {
	ts := runTestServer()

	ts.Repository.CreateDatabase(repositorymodels.Database{ID: testDatabaseName})
	ts.Repository.CreateCollection(testDatabaseName, repositorymodels.Collection{
		ID: testCollectionName,
		PartitionKey: struct {
			Paths   []string "json:\"paths\""
			Kind    string   "json:\"kind\""
			Version int      "json:\"Version\""
		}{
			Paths: []string{"/pk"},
		},
	})
	ts.Repository.CreateDocument(testDatabaseName, testCollectionName, map[string]interface{}{"id": "12345", "pk": "123", "isCool": false, "arr": []int{1, 2, 3}})
	ts.Repository.CreateDocument(testDatabaseName, testCollectionName, map[string]interface{}{"id": "67890", "pk": "456", "isCool": true, "arr": []int{6, 7, 8}})

	client, err := azcosmos.NewClientFromConnectionString(
		fmt.Sprintf("AccountEndpoint=%s;AccountKey=%s", ts.URL, config.DefaultAccountKey),
		&azcosmos.ClientOptions{},
	)
	assert.Nil(t, err)

	collectionClient, err := client.NewContainer(testDatabaseName, testCollectionName)
	assert.Nil(t, err)

	return ts, collectionClient
}

func Test_Documents(t *testing.T) {
	ts, collectionClient := documents_InitializeDb(t)
	defer ts.Server.Close()

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

	t.Run("Should query document with query parameters as accessor", func(t *testing.T) {
		testCosmosQuery(t, collectionClient,
			`select c.id
			FROM c
			WHERE c[@param]="67890"
			ORDER BY c.id`,
			[]azcosmos.QueryParameter{
				{Name: "@param", Value: "id"},
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

	t.Run("Should handle parallel writes", func(t *testing.T) {
		var wg sync.WaitGroup
		rutineCount := 100
		results := make(chan error, rutineCount)

		createCall := func(i int) {
			defer wg.Done()
			item := map[string]interface{}{
				"id":  fmt.Sprintf("id-%d", i),
				"pk":  fmt.Sprintf("pk-%d", i),
				"val": i,
			}
			bytes, err := json.Marshal(item)
			if err != nil {
				results <- err
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			_, err = collectionClient.CreateItem(
				ctx,
				azcosmos.PartitionKey{},
				bytes,
				&azcosmos.ItemOptions{
					EnableContentResponseOnWrite: false,
				},
			)
			results <- err

			collectionClient.ReadItem(ctx, azcosmos.PartitionKey{}, fmt.Sprintf("id-%d", i), nil)
			collectionClient.DeleteItem(ctx, azcosmos.PartitionKey{}, fmt.Sprintf("id-%d", i), nil)
		}

		for i := 0; i < rutineCount; i++ {
			wg.Add(1)
			go createCall(i)
		}

		wg.Wait()
		close(results)

		for err := range results {
			if err != nil {
				t.Errorf("Error creating item: %v", err)
			}
		}
	})
}

func Test_Documents_Patch(t *testing.T) {
	ts, collectionClient := documents_InitializeDb(t)
	defer ts.Server.Close()

	t.Run("Should PATCH document", func(t *testing.T) {
		context := context.TODO()
		expectedData := map[string]interface{}{"id": "67890", "pk": "666", "newField": "newValue", "incr": 15., "setted": "isSet"}

		patch := azcosmos.PatchOperations{}
		patch.AppendAdd("/newField", "newValue")
		patch.AppendIncrement("/incr", 15)
		patch.AppendRemove("/isCool")
		patch.AppendReplace("/pk", "666")
		patch.AppendSet("/setted", "isSet")

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

		var itemResponseBody map[string]interface{}
		json.Unmarshal(itemResponse.Value, &itemResponseBody)

		assert.Equal(t, expectedData["id"], itemResponseBody["id"])
		assert.Equal(t, expectedData["pk"], itemResponseBody["pk"])
		assert.Empty(t, itemResponseBody["isCool"])
		assert.Equal(t, expectedData["newField"], itemResponseBody["newField"])
		assert.Equal(t, expectedData["incr"], itemResponseBody["incr"])
		assert.Equal(t, expectedData["setted"], itemResponseBody["setted"])
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

		r, err := collectionClient.CreateItem(
			context,
			azcosmos.PartitionKey{},
			bytes,
			&azcosmos.ItemOptions{
				EnableContentResponseOnWrite: false,
			},
		)
		assert.NotNil(t, r)
		assert.NotNil(t, err)

		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) {
			assert.Equal(t, http.StatusConflict, respErr.StatusCode)
		} else {
			panic(err)
		}
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

func Test_Documents_TransactionalBatch(t *testing.T) {
	ts, collectionClient := documents_InitializeDb(t)
	defer ts.Server.Close()

	t.Run("Should execute CREATE transactional batch", func(t *testing.T) {
		context := context.TODO()
		batch := collectionClient.NewTransactionalBatch(azcosmos.NewPartitionKeyString("pk"))

		newItem := map[string]interface{}{
			"id": "678901",
		}
		bytes, err := json.Marshal(newItem)
		assert.Nil(t, err)

		batch.CreateItem(bytes, nil)
		response, err := collectionClient.ExecuteTransactionalBatch(context, batch, &azcosmos.TransactionalBatchOptions{})
		assert.Nil(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, 1, len(response.OperationResults))

		operationResponse := response.OperationResults[0]
		assert.NotNil(t, operationResponse)
		assert.NotNil(t, operationResponse.ResourceBody)
		assert.Equal(t, int32(http.StatusCreated), operationResponse.StatusCode)

		var itemResponseBody map[string]interface{}
		json.Unmarshal(operationResponse.ResourceBody, &itemResponseBody)
		assert.Equal(t, newItem["id"], itemResponseBody["id"])

		createdDoc, _ := ts.Repository.GetDocument(testDatabaseName, testCollectionName, newItem["id"].(string))
		assert.Equal(t, newItem["id"], createdDoc["id"])
	})

	t.Run("Should execute DELETE transactional batch", func(t *testing.T) {
		context := context.TODO()
		batch := collectionClient.NewTransactionalBatch(azcosmos.NewPartitionKeyString("pk"))

		batch.DeleteItem("12345", nil)
		response, err := collectionClient.ExecuteTransactionalBatch(context, batch, &azcosmos.TransactionalBatchOptions{})
		assert.Nil(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, 1, len(response.OperationResults))

		operationResponse := response.OperationResults[0]
		assert.NotNil(t, operationResponse)
		assert.Equal(t, int32(http.StatusNoContent), operationResponse.StatusCode)

		_, status := ts.Repository.GetDocument(testDatabaseName, testCollectionName, "12345")
		assert.Equal(t, repositorymodels.StatusNotFound, int(status))
	})

	t.Run("Should execute REPLACE transactional batch", func(t *testing.T) {
		context := context.TODO()
		batch := collectionClient.NewTransactionalBatch(azcosmos.NewPartitionKeyString("pk"))

		newItem := map[string]interface{}{
			"id": "67890",
			"pk": "666",
		}
		bytes, err := json.Marshal(newItem)
		assert.Nil(t, err)

		batch.ReplaceItem("67890", bytes, nil)
		response, err := collectionClient.ExecuteTransactionalBatch(context, batch, &azcosmos.TransactionalBatchOptions{})
		assert.Nil(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, 1, len(response.OperationResults))

		operationResponse := response.OperationResults[0]
		assert.NotNil(t, operationResponse)
		assert.NotNil(t, operationResponse.ResourceBody)
		assert.Equal(t, int32(http.StatusCreated), operationResponse.StatusCode)

		var itemResponseBody map[string]interface{}
		json.Unmarshal(operationResponse.ResourceBody, &itemResponseBody)
		assert.Equal(t, newItem["id"], itemResponseBody["id"])
		assert.Equal(t, newItem["pk"], itemResponseBody["pk"])

		updatedDoc, _ := ts.Repository.GetDocument(testDatabaseName, testCollectionName, newItem["id"].(string))
		assert.Equal(t, newItem["id"], updatedDoc["id"])
		assert.Equal(t, newItem["pk"], updatedDoc["pk"])
	})

	t.Run("Should execute UPSERT transactional batch", func(t *testing.T) {
		context := context.TODO()
		batch := collectionClient.NewTransactionalBatch(azcosmos.NewPartitionKeyString("pk"))

		newItem := map[string]interface{}{
			"id": "678901",
			"pk": "666",
		}
		bytes, err := json.Marshal(newItem)
		assert.Nil(t, err)

		batch.UpsertItem(bytes, nil)
		response, err := collectionClient.ExecuteTransactionalBatch(context, batch, &azcosmos.TransactionalBatchOptions{})
		assert.Nil(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, 1, len(response.OperationResults))

		operationResponse := response.OperationResults[0]
		assert.NotNil(t, operationResponse)
		assert.NotNil(t, operationResponse.ResourceBody)
		assert.Equal(t, int32(http.StatusCreated), operationResponse.StatusCode)

		var itemResponseBody map[string]interface{}
		json.Unmarshal(operationResponse.ResourceBody, &itemResponseBody)
		assert.Equal(t, newItem["id"], itemResponseBody["id"])
		assert.Equal(t, newItem["pk"], itemResponseBody["pk"])

		updatedDoc, _ := ts.Repository.GetDocument(testDatabaseName, testCollectionName, newItem["id"].(string))
		assert.Equal(t, newItem["id"], updatedDoc["id"])
		assert.Equal(t, newItem["pk"], updatedDoc["pk"])
	})

	t.Run("Should execute READ transactional batch", func(t *testing.T) {
		context := context.TODO()
		batch := collectionClient.NewTransactionalBatch(azcosmos.NewPartitionKeyString("pk"))

		batch.ReadItem("67890", nil)
		response, err := collectionClient.ExecuteTransactionalBatch(context, batch, &azcosmos.TransactionalBatchOptions{})
		assert.Nil(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, 1, len(response.OperationResults))

		operationResponse := response.OperationResults[0]
		assert.NotNil(t, operationResponse)
		assert.NotNil(t, operationResponse.ResourceBody)
		assert.Equal(t, int32(http.StatusOK), operationResponse.StatusCode)

		var itemResponseBody map[string]interface{}
		json.Unmarshal(operationResponse.ResourceBody, &itemResponseBody)
		assert.Equal(t, "67890", itemResponseBody["id"])
	})
}
