package badgerdatastore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/logger"
	"github.com/pikami/cosmium/internal/resourceid"
)

func (r *BadgerDataStore) GetAllStoredProcedures(databaseId string, collectionId string) ([]datastore.StoredProcedure, datastore.DataStoreStatus) {
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

	prefix := generateKey(resourceid.ResourceTypeStoredProcedure, databaseId, collectionId, "") + "/"
	storedProcedures, status := listByPrefix[datastore.StoredProcedure](r.db, prefix)
	if status == datastore.StatusOk {
		return storedProcedures, datastore.StatusOk
	}

	return nil, status
}

func (r *BadgerDataStore) GetStoredProcedure(databaseId string, collectionId string, storedProcedureId string) (datastore.StoredProcedure, datastore.DataStoreStatus) {
	storedProcedureKey := generateStoredProcedureKey(databaseId, collectionId, storedProcedureId)

	txn := r.db.NewTransaction(false)
	defer txn.Discard()

	var storedProcedure datastore.StoredProcedure
	status := getKey(txn, storedProcedureKey, &storedProcedure)

	return storedProcedure, status
}

func (r *BadgerDataStore) DeleteStoredProcedure(databaseId string, collectionId string, storedProcedureId string) datastore.DataStoreStatus {
	storedProcedureKey := generateStoredProcedureKey(databaseId, collectionId, storedProcedureId)

	txn := r.db.NewTransaction(true)
	defer txn.Discard()

	exists, err := keyExists(txn, storedProcedureKey)
	if err != nil {
		return datastore.Unknown
	}
	if !exists {
		return datastore.StatusNotFound
	}

	err = txn.Delete([]byte(storedProcedureKey))
	if err != nil {
		logger.ErrorLn("Error while deleting stored procedure:", err)
		return datastore.Unknown
	}

	err = txn.Commit()
	if err != nil {
		logger.ErrorLn("Error while committing transaction:", err)
		return datastore.Unknown
	}

	return datastore.StatusOk
}

func (r *BadgerDataStore) CreateStoredProcedure(databaseId string, collectionId string, storedProcedure datastore.StoredProcedure) (datastore.StoredProcedure, datastore.DataStoreStatus) {
	txn := r.db.NewTransaction(true)
	defer txn.Discard()

	if storedProcedure.ID == "" {
		return datastore.StoredProcedure{}, datastore.BadRequest
	}

	var database datastore.Database
	status := getKey(txn, generateDatabaseKey(databaseId), &database)
	if status != datastore.StatusOk {
		return datastore.StoredProcedure{}, status
	}

	var collection datastore.Collection
	status = getKey(txn, generateCollectionKey(databaseId, collectionId), &collection)
	if status != datastore.StatusOk {
		return datastore.StoredProcedure{}, status
	}

	storedProcedure.TimeStamp = time.Now().Unix()
	storedProcedure.ResourceID = resourceid.NewCombined(collection.ResourceID, resourceid.New(resourceid.ResourceTypeStoredProcedure))
	storedProcedure.ETag = fmt.Sprintf("\"%s\"", uuid.New())
	storedProcedure.Self = fmt.Sprintf("dbs/%s/colls/%s/sprocs/%s/", database.ResourceID, collection.ResourceID, storedProcedure.ResourceID)

	status = insertKey(txn, generateStoredProcedureKey(databaseId, collectionId, storedProcedure.ID), storedProcedure)
	if status != datastore.StatusOk {
		return datastore.StoredProcedure{}, status
	}

	return storedProcedure, datastore.StatusOk
}
