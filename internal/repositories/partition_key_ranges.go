package repositories

import (
	"fmt"

	"github.com/google/uuid"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
	"github.com/pikami/cosmium/internal/resourceid"
)

// I have no idea what this is tbh
func GetPartitionKeyRanges(databaseId string, collectionId string) ([]repositorymodels.PartitionKeyRange, repositorymodels.RepositoryStatus) {
	var ok bool
	var database repositorymodels.Database
	var collection repositorymodels.Collection
	if database, ok = storeState.Databases[databaseId]; !ok {
		return make([]repositorymodels.PartitionKeyRange, 0), repositorymodels.StatusNotFound
	}

	if collection, ok = storeState.Collections[databaseId][collectionId]; !ok {
		return make([]repositorymodels.PartitionKeyRange, 0), repositorymodels.StatusNotFound
	}

	pkrResourceId := resourceid.NewCombined(database.ResourceID, collection.ResourceID, resourceid.New())
	pkrSelf := fmt.Sprintf("dbs/%s/colls/%s/pkranges/%s/", database.ResourceID, collection.ResourceID, pkrResourceId)
	etag := fmt.Sprintf("\"%s\"", uuid.New())

	return []repositorymodels.PartitionKeyRange{
		{
			ResourceID:         pkrResourceId,
			ID:                 "0",
			Etag:               etag,
			MinInclusive:       "",
			MaxExclusive:       "FF",
			RidPrefix:          0,
			Self:               pkrSelf,
			ThroughputFraction: 1,
			Status:             "online",
			Parents:            []interface{}{},
			TimeStamp:          collection.TimeStamp,
			Lsn:                17,
		},
	}, repositorymodels.StatusOk
}
