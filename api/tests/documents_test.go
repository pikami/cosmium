package tests_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
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

func Test_Documents(t *testing.T) {
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
	repositories.CreateDocument(testDatabaseName, testCollectionName, map[string]interface{}{"id": "12345", "pk": "123", "isCool": false})
	repositories.CreateDocument(testDatabaseName, testCollectionName, map[string]interface{}{"id": "67890", "pk": "456", "isCool": true})

	ts := runTestServer()
	defer ts.Close()

	client, err := azcosmos.NewClientFromConnectionString(
		fmt.Sprintf("AccountEndpoint=%s;AccountKey=%s", ts.URL, "asas"),
		&azcosmos.ClientOptions{},
	)
	assert.Nil(t, err)

	collectionClient, err := client.NewContainer(testDatabaseName, testCollectionName)
	assert.Nil(t, err)

	t.Run("Should query document", func(t *testing.T) {
		testCosmosQuery(t, collectionClient,
			"SELECT c.id, c[\"pk\"] FROM c",
			nil,
			[]interface{}{
				map[string]interface{}{"id": "12345", "pk": "123"},
				map[string]interface{}{"id": "67890", "pk": "456"},
			},
		)
	})

	t.Run("Should query VALUE array", func(t *testing.T) {
		testCosmosQuery(t, collectionClient,
			"SELECT VALUE [c.id, c[\"pk\"]] FROM c",
			nil,
			[]interface{}{
				[]interface{}{"12345", "123"},
				[]interface{}{"67890", "456"},
			},
		)
	})

	t.Run("Should query VALUE object", func(t *testing.T) {
		testCosmosQuery(t, collectionClient,
			"SELECT VALUE { id: c.id, _pk: c.pk } FROM c",
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
			WHERE c.isCool=true`,
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
			WHERE c.id=@param_id`,
			[]azcosmos.QueryParameter{
				{Name: "@param_id", Value: "67890"},
			},
			[]interface{}{
				map[string]interface{}{"id": "67890"},
			},
		)
	})
}
