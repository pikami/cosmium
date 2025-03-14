package jsondatastore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/resourceid"
	"golang.org/x/exp/maps"
)

func (r *JsonDataStore) GetAllStoredProcedures(databaseId string, collectionId string) ([]datastore.StoredProcedure, datastore.DataStoreStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	return maps.Values(r.storeState.StoredProcedures[databaseId][collectionId]), datastore.StatusOk
}

func (r *JsonDataStore) GetStoredProcedure(databaseId string, collectionId string, spId string) (datastore.StoredProcedure, datastore.DataStoreStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return datastore.StoredProcedure{}, datastore.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return datastore.StoredProcedure{}, datastore.StatusNotFound
	}

	if sp, ok := r.storeState.StoredProcedures[databaseId][collectionId][spId]; ok {
		return sp, datastore.StatusOk
	}

	return datastore.StoredProcedure{}, datastore.StatusNotFound
}

func (r *JsonDataStore) DeleteStoredProcedure(databaseId string, collectionId string, spId string) datastore.DataStoreStatus {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return datastore.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return datastore.StatusNotFound
	}

	if _, ok := r.storeState.StoredProcedures[databaseId][collectionId][spId]; !ok {
		return datastore.StatusNotFound
	}

	delete(r.storeState.StoredProcedures[databaseId][collectionId], spId)

	return datastore.StatusOk
}

func (r *JsonDataStore) CreateStoredProcedure(databaseId string, collectionId string, sp datastore.StoredProcedure) (datastore.StoredProcedure, datastore.DataStoreStatus) {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	var ok bool
	var database datastore.Database
	var collection datastore.Collection
	if sp.ID == "" {
		return datastore.StoredProcedure{}, datastore.BadRequest
	}

	if database, ok = r.storeState.Databases[databaseId]; !ok {
		return datastore.StoredProcedure{}, datastore.StatusNotFound
	}

	if collection, ok = r.storeState.Collections[databaseId][collectionId]; !ok {
		return datastore.StoredProcedure{}, datastore.StatusNotFound
	}

	if _, ok = r.storeState.StoredProcedures[databaseId][collectionId][sp.ID]; ok {
		return datastore.StoredProcedure{}, datastore.Conflict
	}

	sp.TimeStamp = time.Now().Unix()
	sp.ResourceID = resourceid.NewCombined(collection.ResourceID, resourceid.New(resourceid.ResourceTypeStoredProcedure))
	sp.ETag = fmt.Sprintf("\"%s\"", uuid.New())
	sp.Self = fmt.Sprintf("dbs/%s/colls/%s/sprocs/%s/", database.ResourceID, collection.ResourceID, sp.ResourceID)

	r.storeState.StoredProcedures[databaseId][collectionId][sp.ID] = sp

	return sp, datastore.StatusOk
}
