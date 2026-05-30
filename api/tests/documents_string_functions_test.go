package tests_test

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/stretchr/testify/assert"
)

func documents_InitializeSingleDocumentDb(t *testing.T, ts *TestServer) *azcosmos.ContainerClient {
	ts.DataStore.CreateDatabase(datastore.Database{ID: testDatabaseName})
	ts.DataStore.CreateCollection(testDatabaseName, datastore.Collection{
		ID: testCollectionName,
		PartitionKey: struct {
			Paths   []string "json:\"paths\""
			Kind    string   "json:\"kind\""
			Version int      "json:\"Version\""
		}{
			Paths: []string{"/pk"},
		},
	})
	ts.DataStore.CreateDocument(testDatabaseName, testCollectionName, map[string]interface{}{"id": "regexmatch-test", "pk": "regexmatch-test"})

	client, err := azcosmos.NewClientFromConnectionString(
		fmt.Sprintf("AccountEndpoint=%s;AccountKey=%s", ts.URL, config.DefaultAccountKey),
		&azcosmos.ClientOptions{},
	)
	assert.Nil(t, err)

	collectionClient, err := client.NewContainer(testDatabaseName, testCollectionName)
	assert.Nil(t, err)

	return collectionClient
}

func Test_Documents_RegexMatch(t *testing.T) {
	presets := []testPreset{PresetJsonStore, PresetBadgerStore}

	runTestsWithPresets(t, "Test_Documents_RegexMatch", presets, func(t *testing.T, ts *TestServer, client *azcosmos.Client) {
		collectionClient := documents_InitializeSingleDocumentDb(t, ts)

		t.Run("Should execute REGEXMATCH()", func(t *testing.T) {
			testCosmosQuery(t, collectionClient,
				`SELECT VALUE {
					noModifiers: REGEXMATCH("abcd", "ABC"),
					caseInsensitive: REGEXMATCH("abcd", "ABC", "i"),
					wildcardCharacter: REGEXMATCH("abcd", "ab.", ""),
					ignoreWhiteSpace: REGEXMATCH("abcd", "ab c", "x"),
					caseInsensitiveAndIgnoreWhiteSpace: REGEXMATCH("abcd", "aB c", "ix"),
					containNumberBetweenZeroAndNine: REGEXMATCH("03a", "[0-9]"),
					containPrefix: REGEXMATCH("salt3824908", "salt{1}"),
					containsFiveLetterWordStartingWithS: REGEXMATCH("shame", "s....", "i")
				}`,
				nil,
				[]interface{}{
					map[string]interface{}{
						"noModifiers":                         false,
						"caseInsensitive":                     true,
						"wildcardCharacter":                   true,
						"ignoreWhiteSpace":                    true,
						"caseInsensitiveAndIgnoreWhiteSpace":  true,
						"containNumberBetweenZeroAndNine":     true,
						"containPrefix":                       true,
						"containsFiveLetterWordStartingWithS": true,
					},
				},
			)
		})
	})
}
