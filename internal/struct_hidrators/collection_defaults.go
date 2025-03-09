package structhidrators

import "github.com/pikami/cosmium/internal/datastore"

var defaultCollection datastore.Collection = datastore.Collection{
	IndexingPolicy: datastore.CollectionIndexingPolicy{
		IndexingMode: "consistent",
		Automatic:    true,
		IncludedPaths: []datastore.CollectionIndexingPolicyPath{
			{Path: "/*"},
		},
		ExcludedPaths: []datastore.CollectionIndexingPolicyPath{
			{Path: "/\"_etag\"/?"},
		},
	},
	PartitionKey: datastore.CollectionPartitionKey{
		Paths:   []string{"/_partitionKey"},
		Kind:    "Hash",
		Version: 2,
	},
	ResourceID: "nFFFFFFFFFF=",
	TimeStamp:  0,
	Self:       "",
	ETag:       "\"00000000-0000-0000-0000-000000000000\"",
	Docs:       "docs/",
	Sprocs:     "sprocs/",
	Triggers:   "triggers/",
	Udfs:       "udfs/",
	Conflicts:  "conflicts/",
}
