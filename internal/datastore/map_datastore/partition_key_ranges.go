package mapdatastore

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/resourceid"
)

// I have no idea what this is tbh
func (r *MapDataStore) GetPartitionKeyRanges(databaseId string, collectionId string) ([]datastore.PartitionKeyRange, datastore.DataStoreStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	databaseRid := databaseId
	collectionRid := collectionId
	var timestamp int64 = 0

	if database, ok := r.storeState.Databases[databaseId]; !ok {
		databaseRid = database.ResourceID
	}

	if collection, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		collectionRid = collection.ResourceID
		timestamp = collection.TimeStamp
	}

	pkrResourceId := resourceid.NewCombined(collectionRid, resourceid.New(resourceid.ResourceTypePartitionKeyRange))
	pkrSelf := fmt.Sprintf("dbs/%s/colls/%s/pkranges/%s/", databaseRid, collectionRid, pkrResourceId)
	etag := fmt.Sprintf("\"%s\"", uuid.New())

	return []datastore.PartitionKeyRange{
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
			TimeStamp:          timestamp,
			Lsn:                17,
		},
	}, datastore.StatusOk
}
