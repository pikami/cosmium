package badgerdatastore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/logger"
	"github.com/pikami/cosmium/internal/resourceid"
)

func (r *BadgerDataStore) GetAllTriggers(databaseId string, collectionId string) ([]datastore.Trigger, datastore.DataStoreStatus) {
	txn := r.db.NewTransaction(false)
	defer txn.Discard()

	dbExists, err := keyExists(txn, generateDatabaseKey(databaseId))
	if err != nil || !dbExists {
		return nil, datastore.StatusNotFound
	}

	collExists, err := keyExists(txn, generateCollectionKey(databaseId, collectionId))
	if err != nil || !collExists {
		return nil, datastore.StatusNotFound
	}

	triggers, status := listByPrefix[datastore.Trigger](r.db, generateKey(resourceid.ResourceTypeTrigger, databaseId, collectionId, ""))
	if status == datastore.StatusOk {
		return triggers, datastore.StatusOk
	}

	return nil, status
}

func (r *BadgerDataStore) GetTrigger(databaseId string, collectionId string, triggerId string) (datastore.Trigger, datastore.DataStoreStatus) {
	triggerKey := generateTriggerKey(databaseId, collectionId, triggerId)

	txn := r.db.NewTransaction(false)
	defer txn.Discard()

	var trigger datastore.Trigger
	status := getKey(txn, triggerKey, &trigger)

	return trigger, status
}

func (r *BadgerDataStore) DeleteTrigger(databaseId string, collectionId string, triggerId string) datastore.DataStoreStatus {
	triggerKey := generateTriggerKey(databaseId, collectionId, triggerId)

	txn := r.db.NewTransaction(true)
	defer txn.Discard()

	exists, err := keyExists(txn, triggerKey)
	if err != nil {
		return datastore.Unknown
	}
	if !exists {
		return datastore.StatusNotFound
	}

	err = txn.Delete([]byte(triggerKey))
	if err != nil {
		logger.ErrorLn("Error while deleting trigger:", err)
		return datastore.Unknown
	}

	err = txn.Commit()
	if err != nil {
		logger.ErrorLn("Error while committing transaction:", err)
		return datastore.Unknown
	}

	return datastore.StatusOk
}

func (r *BadgerDataStore) CreateTrigger(databaseId string, collectionId string, trigger datastore.Trigger) (datastore.Trigger, datastore.DataStoreStatus) {
	txn := r.db.NewTransaction(true)
	defer txn.Discard()

	if trigger.ID == "" {
		return datastore.Trigger{}, datastore.BadRequest
	}

	var database datastore.Database
	status := getKey(txn, generateDatabaseKey(databaseId), &database)
	if status != datastore.StatusOk {
		return datastore.Trigger{}, status
	}

	var collection datastore.Collection
	status = getKey(txn, generateCollectionKey(databaseId, collectionId), &collection)
	if status != datastore.StatusOk {
		return datastore.Trigger{}, status
	}

	trigger.TimeStamp = time.Now().Unix()
	trigger.ResourceID = resourceid.NewCombined(collection.ResourceID, resourceid.New(resourceid.ResourceTypeTrigger))
	trigger.ETag = fmt.Sprintf("\"%s\"", uuid.New())
	trigger.Self = fmt.Sprintf("dbs/%s/colls/%s/triggers/%s/", database.ResourceID, collection.ResourceID, trigger.ResourceID)

	status = insertKey(txn, generateTriggerKey(databaseId, collectionId, trigger.ID), trigger)
	if status != datastore.StatusOk {
		return datastore.Trigger{}, status
	}

	return trigger, datastore.StatusOk
}
