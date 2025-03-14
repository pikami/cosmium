package jsondatastore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/resourceid"
	"golang.org/x/exp/maps"
)

func (r *JsonDataStore) GetAllTriggers(databaseId string, collectionId string) ([]datastore.Trigger, datastore.DataStoreStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	return maps.Values(r.storeState.Triggers[databaseId][collectionId]), datastore.StatusOk
}

func (r *JsonDataStore) GetTrigger(databaseId string, collectionId string, triggerId string) (datastore.Trigger, datastore.DataStoreStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return datastore.Trigger{}, datastore.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return datastore.Trigger{}, datastore.StatusNotFound
	}

	if trigger, ok := r.storeState.Triggers[databaseId][collectionId][triggerId]; ok {
		return trigger, datastore.StatusOk
	}

	return datastore.Trigger{}, datastore.StatusNotFound
}

func (r *JsonDataStore) DeleteTrigger(databaseId string, collectionId string, triggerId string) datastore.DataStoreStatus {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return datastore.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return datastore.StatusNotFound
	}

	if _, ok := r.storeState.Triggers[databaseId][collectionId][triggerId]; !ok {
		return datastore.StatusNotFound
	}

	delete(r.storeState.Triggers[databaseId][collectionId], triggerId)

	return datastore.StatusOk
}

func (r *JsonDataStore) CreateTrigger(databaseId string, collectionId string, trigger datastore.Trigger) (datastore.Trigger, datastore.DataStoreStatus) {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	var ok bool
	var database datastore.Database
	var collection datastore.Collection
	if trigger.ID == "" {
		return datastore.Trigger{}, datastore.BadRequest
	}

	if database, ok = r.storeState.Databases[databaseId]; !ok {
		return datastore.Trigger{}, datastore.StatusNotFound
	}

	if collection, ok = r.storeState.Collections[databaseId][collectionId]; !ok {
		return datastore.Trigger{}, datastore.StatusNotFound
	}

	if _, ok = r.storeState.Triggers[databaseId][collectionId][trigger.ID]; ok {
		return datastore.Trigger{}, datastore.Conflict
	}

	trigger.TimeStamp = time.Now().Unix()
	trigger.ResourceID = resourceid.NewCombined(collection.ResourceID, resourceid.New(resourceid.ResourceTypeTrigger))
	trigger.ETag = fmt.Sprintf("\"%s\"", uuid.New())
	trigger.Self = fmt.Sprintf("dbs/%s/colls/%s/triggers/%s/", database.ResourceID, collection.ResourceID, trigger.ResourceID)

	r.storeState.Triggers[databaseId][collectionId][trigger.ID] = trigger

	return trigger, datastore.StatusOk
}
