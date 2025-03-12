package badgerdatastore

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/resourceid"
)

// I have no idea what this is tbh
func (r *BadgerDataStore) GetPartitionKeyRanges(databaseId string, collectionId string) ([]datastore.PartitionKeyRange, datastore.DataStoreStatus) {
	databaseRid := databaseId
	collectionRid := collectionId
	var timestamp int64 = 0

	txn := r.db.NewTransaction(false)
	defer txn.Discard()

	var database datastore.Database
	status := getKey(txn, generateDatabaseKey(databaseId), &database)
	if status != datastore.StatusOk {
		databaseRid = database.ResourceID
	}

	var collection datastore.Collection
	status = getKey(txn, generateCollectionKey(databaseId, collectionId), &collection)
	if status != datastore.StatusOk {
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
