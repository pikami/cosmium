package tests_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func Test_Documents_ArrayContains(t *testing.T) {
	presets := []testPreset{PresetJsonStore, PresetBadgerStore}

	runTestsWithPresets(t, "Test_Documents_ArrayContains", presets, func(t *testing.T, ts *TestServer, client *azcosmos.Client) {
		collectionClient := documents_InitializeDb(t, ts)

		t.Run("Should execute ARRAY_CONTAINS() without partial match argument", func(t *testing.T) {
			testCosmosQuery(t, collectionClient,
				`SELECT VALUE ARRAY_CONTAINS(["apple", "banana", "cherry"], "banana") FROM c ORDER BY c.id`,
				nil,
				[]interface{}{true, true},
			)
		})

		t.Run("Should execute ARRAY_CONTAINS() returning false for missing item", func(t *testing.T) {
			testCosmosQuery(t, collectionClient,
				`SELECT VALUE ARRAY_CONTAINS(["apple", "banana", "cherry"], "grape") FROM c ORDER BY c.id`,
				nil,
				[]interface{}{false, false},
			)
		})

		t.Run("Should execute ARRAY_CONTAINS() with object full match", func(t *testing.T) {
			testCosmosQuery(t, collectionClient,
				`SELECT VALUE ARRAY_CONTAINS([{"name": "apple", "color": "red"}], {"name": "apple"}) FROM c ORDER BY c.id`,
				nil,
				[]interface{}{false, false},
			)
		})

		t.Run("Should execute ARRAY_CONTAINS() with object partial match", func(t *testing.T) {
			testCosmosQuery(t, collectionClient,
				`SELECT VALUE ARRAY_CONTAINS([{"name": "apple", "color": "red"}], {"name": "apple"}, true) FROM c ORDER BY c.id`,
				nil,
				[]interface{}{true, true},
			)
		})
	})
}
