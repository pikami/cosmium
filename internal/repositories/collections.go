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

func GetAllCollections(databaseId string) ([]repositorymodels.Collection, repositorymodels.RepositoryStatus) {
	storeState.RLock()
	defer storeState.RUnlock()

	if _, ok := storeState.Databases[databaseId]; !ok {
		return make([]repositorymodels.Collection, 0), repositorymodels.StatusNotFound
	}

	return maps.Values(storeState.Collections[databaseId]), repositorymodels.StatusOk
}

func GetCollection(databaseId string, collectionId string) (repositorymodels.Collection, repositorymodels.RepositoryStatus) {
	storeState.RLock()
	defer storeState.RUnlock()

	if _, ok := storeState.Databases[databaseId]; !ok {
		return repositorymodels.Collection{}, repositorymodels.StatusNotFound
	}

	if _, ok := storeState.Collections[databaseId][collectionId]; !ok {
		return repositorymodels.Collection{}, repositorymodels.StatusNotFound
	}

	return storeState.Collections[databaseId][collectionId], repositorymodels.StatusOk
}

func DeleteCollection(databaseId string, collectionId string) repositorymodels.RepositoryStatus {
	storeState.Lock()
	defer storeState.Unlock()

	if _, ok := storeState.Databases[databaseId]; !ok {
		return repositorymodels.StatusNotFound
	}

	if _, ok := storeState.Collections[databaseId][collectionId]; !ok {
		return repositorymodels.StatusNotFound
	}

	delete(storeState.Collections[databaseId], collectionId)

	return repositorymodels.StatusOk
}

func CreateCollection(databaseId string, newCollection repositorymodels.Collection) (repositorymodels.Collection, repositorymodels.RepositoryStatus) {
	storeState.Lock()
	defer storeState.Unlock()

	var ok bool
	var database repositorymodels.Database
	if database, ok = storeState.Databases[databaseId]; !ok {
		return repositorymodels.Collection{}, repositorymodels.StatusNotFound
	}

	if _, ok = storeState.Collections[databaseId][newCollection.ID]; ok {
		return repositorymodels.Collection{}, repositorymodels.Conflict
	}

	newCollection = structhidrators.Hidrate(newCollection).(repositorymodels.Collection)

	newCollection.TimeStamp = time.Now().Unix()
	newCollection.ResourceID = resourceid.NewCombined(database.ResourceID, resourceid.New())
	newCollection.ETag = fmt.Sprintf("\"%s\"", uuid.New())
	newCollection.Self = fmt.Sprintf("dbs/%s/colls/%s/", database.ResourceID, newCollection.ResourceID)

	storeState.Collections[databaseId][newCollection.ID] = newCollection
	storeState.Documents[databaseId][newCollection.ID] = make(map[string]repositorymodels.Document)

	return newCollection, repositorymodels.StatusOk
}
