package structhidrators

import (
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

var defaultCollection repositorymodels.Collection = repositorymodels.Collection{
	IndexingPolicy: repositorymodels.CollectionIndexingPolicy{
		IndexingMode: "consistent",
		Automatic:    true,
		IncludedPaths: []repositorymodels.CollectionIndexingPolicyPath{
			{Path: "/*"},
		},
		ExcludedPaths: []repositorymodels.CollectionIndexingPolicyPath{
			{Path: "/\"_etag\"/?"},
		},
	},
	PartitionKey: repositorymodels.CollectionPartitionKey{
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
