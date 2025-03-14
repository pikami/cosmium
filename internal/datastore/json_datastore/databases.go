package jsondatastore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/resourceid"
	"golang.org/x/exp/maps"
)

func (r *JsonDataStore) GetAllDatabases() ([]datastore.Database, datastore.DataStoreStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	return maps.Values(r.storeState.Databases), datastore.StatusOk
}

func (r *JsonDataStore) GetDatabase(id string) (datastore.Database, datastore.DataStoreStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	if database, ok := r.storeState.Databases[id]; ok {
		return database, datastore.StatusOk
	}

	return datastore.Database{}, datastore.StatusNotFound
}

func (r *JsonDataStore) DeleteDatabase(id string) datastore.DataStoreStatus {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	if _, ok := r.storeState.Databases[id]; !ok {
		return datastore.StatusNotFound
	}

	delete(r.storeState.Databases, id)
	delete(r.storeState.Collections, id)
	delete(r.storeState.Documents, id)
	delete(r.storeState.Triggers, id)
	delete(r.storeState.StoredProcedures, id)
	delete(r.storeState.UserDefinedFunctions, id)

	return datastore.StatusOk
}

func (r *JsonDataStore) CreateDatabase(newDatabase datastore.Database) (datastore.Database, datastore.DataStoreStatus) {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	if _, ok := r.storeState.Databases[newDatabase.ID]; ok {
		return datastore.Database{}, datastore.Conflict
	}

	newDatabase.TimeStamp = time.Now().Unix()
	newDatabase.ResourceID = resourceid.New(resourceid.ResourceTypeDatabase)
	newDatabase.ETag = fmt.Sprintf("\"%s\"", uuid.New())
	newDatabase.Self = fmt.Sprintf("dbs/%s/", newDatabase.ResourceID)

	r.storeState.Databases[newDatabase.ID] = newDatabase
	r.storeState.Collections[newDatabase.ID] = make(map[string]datastore.Collection)
	r.storeState.Documents[newDatabase.ID] = make(map[string]map[string]datastore.Document)
	r.storeState.Triggers[newDatabase.ID] = make(map[string]map[string]datastore.Trigger)
	r.storeState.StoredProcedures[newDatabase.ID] = make(map[string]map[string]datastore.StoredProcedure)
	r.storeState.UserDefinedFunctions[newDatabase.ID] = make(map[string]map[string]datastore.UserDefinedFunction)

	return newDatabase, datastore.StatusOk
}
