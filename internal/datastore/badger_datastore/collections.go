package badgerdatastore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/logger"
	"github.com/pikami/cosmium/internal/resourceid"
	structhidrators "github.com/pikami/cosmium/internal/struct_hidrators"
)

func (r *BadgerDataStore) GetAllCollections(databaseId string) ([]datastore.Collection, datastore.DataStoreStatus) {
	exists, err := keyExists(r.db.NewTransaction(false), generateDatabaseKey(databaseId))
	if err != nil {
		logger.ErrorLn("Error while checking if database exists:", err)
		return nil, datastore.Unknown
	}

	if !exists {
		return nil, datastore.StatusNotFound
	}

	colls, status := listByPrefix[datastore.Collection](r.db, generateKey(resourceid.ResourceTypeCollection, databaseId, "", ""))
	if status == datastore.StatusOk {
		return colls, datastore.StatusOk
	}

	return nil, status
}

func (r *BadgerDataStore) GetCollection(databaseId string, collectionId string) (datastore.Collection, datastore.DataStoreStatus) {
	collectionKey := generateCollectionKey(databaseId, collectionId)

	txn := r.db.NewTransaction(false)
	defer txn.Discard()

	var collection datastore.Collection
	status := getKey(txn, collectionKey, &collection)

	return collection, status
}

func (r *BadgerDataStore) DeleteCollection(databaseId string, collectionId string) datastore.DataStoreStatus {
	collectionKey := generateCollectionKey(databaseId, collectionId)

	txn := r.db.NewTransaction(true)
	defer txn.Discard()

	prefixes := []string{
		generateKey(resourceid.ResourceTypeDocument, databaseId, collectionId, ""),
		generateKey(resourceid.ResourceTypeTrigger, databaseId, collectionId, ""),
		generateKey(resourceid.ResourceTypeStoredProcedure, databaseId, collectionId, ""),
		generateKey(resourceid.ResourceTypeUserDefinedFunction, databaseId, collectionId, ""),
		collectionKey,
	}
	for _, prefix := range prefixes {
		if err := deleteKeysByPrefix(txn, prefix); err != nil {
			return datastore.Unknown
		}
	}

	err := txn.Commit()
	if err != nil {
		logger.ErrorLn("Error while committing transaction:", err)
		return datastore.Unknown
	}

	return datastore.StatusOk
}

func (r *BadgerDataStore) CreateCollection(databaseId string, newCollection datastore.Collection) (datastore.Collection, datastore.DataStoreStatus) {
	collectionKey := generateCollectionKey(databaseId, newCollection.ID)

	txn := r.db.NewTransaction(true)
	defer txn.Discard()

	collectionExists, err := keyExists(txn, collectionKey)
	if err != nil || collectionExists {
		return datastore.Collection{}, datastore.Conflict
	}

	var database datastore.Database
	status := getKey(txn, generateDatabaseKey(databaseId), &database)
	if status != datastore.StatusOk {
		return datastore.Collection{}, status
	}

	newCollection = structhidrators.Hidrate(newCollection).(datastore.Collection)

	newCollection.TimeStamp = time.Now().Unix()
	newCollection.ResourceID = resourceid.NewCombined(database.ResourceID, resourceid.New(resourceid.ResourceTypeCollection))
	newCollection.ETag = fmt.Sprintf("\"%s\"", uuid.New())
	newCollection.Self = fmt.Sprintf("dbs/%s/colls/%s/", database.ResourceID, newCollection.ResourceID)

	status = insertKey(txn, collectionKey, newCollection)
	if status != datastore.StatusOk {
		return datastore.Collection{}, status
	}

	return newCollection, datastore.StatusOk
}
