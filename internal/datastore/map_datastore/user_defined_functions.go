package mapdatastore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/resourceid"
	"golang.org/x/exp/maps"
)

func (r *MapDataStore) GetAllUserDefinedFunctions(databaseId string, collectionId string) ([]datastore.UserDefinedFunction, datastore.DataStoreStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	return maps.Values(r.storeState.UserDefinedFunctions[databaseId][collectionId]), datastore.StatusOk
}

func (r *MapDataStore) GetUserDefinedFunction(databaseId string, collectionId string, udfId string) (datastore.UserDefinedFunction, datastore.DataStoreStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return datastore.UserDefinedFunction{}, datastore.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return datastore.UserDefinedFunction{}, datastore.StatusNotFound
	}

	if udf, ok := r.storeState.UserDefinedFunctions[databaseId][collectionId][udfId]; ok {
		return udf, datastore.StatusOk
	}

	return datastore.UserDefinedFunction{}, datastore.StatusNotFound
}

func (r *MapDataStore) DeleteUserDefinedFunction(databaseId string, collectionId string, udfId string) datastore.DataStoreStatus {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return datastore.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return datastore.StatusNotFound
	}

	if _, ok := r.storeState.UserDefinedFunctions[databaseId][collectionId][udfId]; !ok {
		return datastore.StatusNotFound
	}

	delete(r.storeState.UserDefinedFunctions[databaseId][collectionId], udfId)

	return datastore.StatusOk
}

func (r *MapDataStore) CreateUserDefinedFunction(databaseId string, collectionId string, udf datastore.UserDefinedFunction) (datastore.UserDefinedFunction, datastore.DataStoreStatus) {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	var ok bool
	var database datastore.Database
	var collection datastore.Collection
	if udf.ID == "" {
		return datastore.UserDefinedFunction{}, datastore.BadRequest
	}

	if database, ok = r.storeState.Databases[databaseId]; !ok {
		return datastore.UserDefinedFunction{}, datastore.StatusNotFound
	}

	if collection, ok = r.storeState.Collections[databaseId][collectionId]; !ok {
		return datastore.UserDefinedFunction{}, datastore.StatusNotFound
	}

	if _, ok := r.storeState.UserDefinedFunctions[databaseId][collectionId][udf.ID]; ok {
		return datastore.UserDefinedFunction{}, datastore.Conflict
	}

	udf.TimeStamp = time.Now().Unix()
	udf.ResourceID = resourceid.NewCombined(collection.ResourceID, resourceid.New(resourceid.ResourceTypeUserDefinedFunction))
	udf.ETag = fmt.Sprintf("\"%s\"", uuid.New())
	udf.Self = fmt.Sprintf("dbs/%s/colls/%s/udfs/%s/", database.ResourceID, collection.ResourceID, udf.ResourceID)

	r.storeState.UserDefinedFunctions[databaseId][collectionId][udf.ID] = udf

	return udf, datastore.StatusOk
}
