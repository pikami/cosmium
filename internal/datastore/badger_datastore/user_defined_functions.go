package badgerdatastore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/logger"
	"github.com/pikami/cosmium/internal/resourceid"
)

func (r *BadgerDataStore) GetAllUserDefinedFunctions(databaseId string, collectionId string) ([]datastore.UserDefinedFunction, datastore.DataStoreStatus) {
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

	prefix := generateKey(resourceid.ResourceTypeUserDefinedFunction, databaseId, collectionId, "") + "/"
	udfs, status := listByPrefix[datastore.UserDefinedFunction](r.db, prefix)
	if status == datastore.StatusOk {
		return udfs, datastore.StatusOk
	}

	return nil, status
}

func (r *BadgerDataStore) GetUserDefinedFunction(databaseId string, collectionId string, udfId string) (datastore.UserDefinedFunction, datastore.DataStoreStatus) {
	udfKey := generateUserDefinedFunctionKey(databaseId, collectionId, udfId)

	txn := r.db.NewTransaction(false)
	defer txn.Discard()

	var udf datastore.UserDefinedFunction
	status := getKey(txn, udfKey, &udf)

	return udf, status
}

func (r *BadgerDataStore) DeleteUserDefinedFunction(databaseId string, collectionId string, udfId string) datastore.DataStoreStatus {
	udfKey := generateUserDefinedFunctionKey(databaseId, collectionId, udfId)

	txn := r.db.NewTransaction(true)
	defer txn.Discard()

	exists, err := keyExists(txn, udfKey)
	if err != nil {
		return datastore.Unknown
	}
	if !exists {
		return datastore.StatusNotFound
	}

	err = txn.Delete([]byte(udfKey))
	if err != nil {
		logger.ErrorLn("Error while deleting user defined function:", err)
		return datastore.Unknown
	}

	err = txn.Commit()
	if err != nil {
		logger.ErrorLn("Error while committing transaction:", err)
		return datastore.Unknown
	}

	return datastore.StatusOk
}

func (r *BadgerDataStore) CreateUserDefinedFunction(databaseId string, collectionId string, udf datastore.UserDefinedFunction) (datastore.UserDefinedFunction, datastore.DataStoreStatus) {
	txn := r.db.NewTransaction(true)
	defer txn.Discard()

	if udf.ID == "" {
		return datastore.UserDefinedFunction{}, datastore.BadRequest
	}

	var database datastore.Database
	status := getKey(txn, generateDatabaseKey(databaseId), &database)
	if status != datastore.StatusOk {
		return datastore.UserDefinedFunction{}, status
	}

	var collection datastore.Collection
	status = getKey(txn, generateCollectionKey(databaseId, collectionId), &collection)
	if status != datastore.StatusOk {
		return datastore.UserDefinedFunction{}, status
	}

	udf.TimeStamp = time.Now().Unix()
	udf.ResourceID = resourceid.NewCombined(collection.ResourceID, resourceid.New(resourceid.ResourceTypeUserDefinedFunction))
	udf.ETag = fmt.Sprintf("\"%s\"", uuid.New())
	udf.Self = fmt.Sprintf("dbs/%s/colls/%s/udfs/%s/", database.ResourceID, collection.ResourceID, udf.ResourceID)

	status = insertKey(txn, generateUserDefinedFunctionKey(databaseId, collectionId, udf.ID), udf)
	if status != datastore.StatusOk {
		return datastore.UserDefinedFunction{}, status
	}

	return udf, datastore.StatusOk
}
