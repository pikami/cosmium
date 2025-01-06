package repositories

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
	"github.com/pikami/cosmium/internal/resourceid"
	"golang.org/x/exp/maps"
)

func (r *DataRepository) GetAllTriggers(databaseId string, collectionId string) ([]repositorymodels.Trigger, repositorymodels.RepositoryStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	return maps.Values(r.storeState.Triggers[databaseId][collectionId]), repositorymodels.StatusOk
}

func (r *DataRepository) GetTrigger(databaseId string, collectionId string, triggerId string) (repositorymodels.Trigger, repositorymodels.RepositoryStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return repositorymodels.Trigger{}, repositorymodels.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return repositorymodels.Trigger{}, repositorymodels.StatusNotFound
	}

	if trigger, ok := r.storeState.Triggers[databaseId][collectionId][triggerId]; ok {
		return trigger, repositorymodels.StatusOk
	}

	return repositorymodels.Trigger{}, repositorymodels.StatusNotFound
}

func (r *DataRepository) DeleteTrigger(databaseId string, collectionId string, triggerId string) repositorymodels.RepositoryStatus {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return repositorymodels.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return repositorymodels.StatusNotFound
	}

	if _, ok := r.storeState.Triggers[databaseId][collectionId][triggerId]; !ok {
		return repositorymodels.StatusNotFound
	}

	delete(r.storeState.Triggers[databaseId][collectionId], triggerId)

	return repositorymodels.StatusOk
}

func (r *DataRepository) CreateTrigger(databaseId string, collectionId string, trigger repositorymodels.Trigger) (repositorymodels.Trigger, repositorymodels.RepositoryStatus) {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	var ok bool
	var database repositorymodels.Database
	var collection repositorymodels.Collection
	if trigger.ID == "" {
		return repositorymodels.Trigger{}, repositorymodels.BadRequest
	}

	if database, ok = r.storeState.Databases[databaseId]; !ok {
		return repositorymodels.Trigger{}, repositorymodels.StatusNotFound
	}

	if collection, ok = r.storeState.Collections[databaseId][collectionId]; !ok {
		return repositorymodels.Trigger{}, repositorymodels.StatusNotFound
	}

	if _, ok = r.storeState.Triggers[databaseId][collectionId][trigger.ID]; ok {
		return repositorymodels.Trigger{}, repositorymodels.Conflict
	}

	trigger.TimeStamp = time.Now().Unix()
	trigger.ResourceID = resourceid.NewCombined(database.ResourceID, collection.ResourceID, resourceid.New())
	trigger.ETag = fmt.Sprintf("\"%s\"", uuid.New())
	trigger.Self = fmt.Sprintf("dbs/%s/colls/%s/triggers/%s/", database.ResourceID, collection.ResourceID, trigger.ResourceID)

	r.storeState.Triggers[databaseId][collectionId][trigger.ID] = trigger

	return trigger, repositorymodels.StatusOk
}
