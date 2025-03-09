package mapdatastore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/resourceid"
	structhidrators "github.com/pikami/cosmium/internal/struct_hidrators"
	"golang.org/x/exp/maps"
)

func (r *MapDataStore) GetAllCollections(databaseId string) ([]datastore.Collection, datastore.DataStoreStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return make([]datastore.Collection, 0), datastore.StatusNotFound
	}

	return maps.Values(r.storeState.Collections[databaseId]), datastore.StatusOk
}

func (r *MapDataStore) GetCollection(databaseId string, collectionId string) (datastore.Collection, datastore.DataStoreStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return datastore.Collection{}, datastore.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return datastore.Collection{}, datastore.StatusNotFound
	}

	return r.storeState.Collections[databaseId][collectionId], datastore.StatusOk
}

func (r *MapDataStore) DeleteCollection(databaseId string, collectionId string) datastore.DataStoreStatus {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return datastore.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return datastore.StatusNotFound
	}

	delete(r.storeState.Collections[databaseId], collectionId)
	delete(r.storeState.Documents[databaseId], collectionId)
	delete(r.storeState.Triggers[databaseId], collectionId)
	delete(r.storeState.StoredProcedures[databaseId], collectionId)
	delete(r.storeState.UserDefinedFunctions[databaseId], collectionId)

	return datastore.StatusOk
}

func (r *MapDataStore) CreateCollection(databaseId string, newCollection datastore.Collection) (datastore.Collection, datastore.DataStoreStatus) {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	var ok bool
	var database datastore.Database
	if database, ok = r.storeState.Databases[databaseId]; !ok {
		return datastore.Collection{}, datastore.StatusNotFound
	}

	if _, ok = r.storeState.Collections[databaseId][newCollection.ID]; ok {
		return datastore.Collection{}, datastore.Conflict
	}

	newCollection = structhidrators.Hidrate(newCollection).(datastore.Collection)

	newCollection.TimeStamp = time.Now().Unix()
	newCollection.ResourceID = resourceid.NewCombined(database.ResourceID, resourceid.New(resourceid.ResourceTypeCollection))
	newCollection.ETag = fmt.Sprintf("\"%s\"", uuid.New())
	newCollection.Self = fmt.Sprintf("dbs/%s/colls/%s/", database.ResourceID, newCollection.ResourceID)

	r.storeState.Collections[databaseId][newCollection.ID] = newCollection
	r.storeState.Documents[databaseId][newCollection.ID] = make(map[string]datastore.Document)
	r.storeState.Triggers[databaseId][newCollection.ID] = make(map[string]datastore.Trigger)
	r.storeState.StoredProcedures[databaseId][newCollection.ID] = make(map[string]datastore.StoredProcedure)
	r.storeState.UserDefinedFunctions[databaseId][newCollection.ID] = make(map[string]datastore.UserDefinedFunction)

	return newCollection, datastore.StatusOk
}
