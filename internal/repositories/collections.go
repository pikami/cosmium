package repositories

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
	"github.com/pikami/cosmium/internal/resourceid"
	structhidrators "github.com/pikami/cosmium/internal/struct_hidrators"
	"golang.org/x/exp/maps"
)

func (r *DataRepository) GetAllCollections(databaseId string) ([]repositorymodels.Collection, repositorymodels.RepositoryStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return make([]repositorymodels.Collection, 0), repositorymodels.StatusNotFound
	}

	return maps.Values(r.storeState.Collections[databaseId]), repositorymodels.StatusOk
}

func (r *DataRepository) GetCollection(databaseId string, collectionId string) (repositorymodels.Collection, repositorymodels.RepositoryStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return repositorymodels.Collection{}, repositorymodels.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return repositorymodels.Collection{}, repositorymodels.StatusNotFound
	}

	return r.storeState.Collections[databaseId][collectionId], repositorymodels.StatusOk
}

func (r *DataRepository) DeleteCollection(databaseId string, collectionId string) repositorymodels.RepositoryStatus {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return repositorymodels.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return repositorymodels.StatusNotFound
	}

	delete(r.storeState.Collections[databaseId], collectionId)

	return repositorymodels.StatusOk
}

func (r *DataRepository) CreateCollection(databaseId string, newCollection repositorymodels.Collection) (repositorymodels.Collection, repositorymodels.RepositoryStatus) {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	var ok bool
	var database repositorymodels.Database
	if database, ok = r.storeState.Databases[databaseId]; !ok {
		return repositorymodels.Collection{}, repositorymodels.StatusNotFound
	}

	if _, ok = r.storeState.Collections[databaseId][newCollection.ID]; ok {
		return repositorymodels.Collection{}, repositorymodels.Conflict
	}

	newCollection = structhidrators.Hidrate(newCollection).(repositorymodels.Collection)

	newCollection.TimeStamp = time.Now().Unix()
	newCollection.ResourceID = resourceid.NewCombined(database.ResourceID, resourceid.New())
	newCollection.ETag = fmt.Sprintf("\"%s\"", uuid.New())
	newCollection.Self = fmt.Sprintf("dbs/%s/colls/%s/", database.ResourceID, newCollection.ResourceID)

	r.storeState.Collections[databaseId][newCollection.ID] = newCollection
	r.storeState.Documents[databaseId][newCollection.ID] = make(map[string]repositorymodels.Document)
	r.storeState.Triggers[databaseId][newCollection.ID] = make(map[string]repositorymodels.Trigger)
	r.storeState.StoredProcedures[databaseId][newCollection.ID] = make(map[string]repositorymodels.StoredProcedure)
	r.storeState.UserDefinedFunctions[databaseId][newCollection.ID] = make(map[string]repositorymodels.UserDefinedFunction)

	return newCollection, repositorymodels.StatusOk
}
