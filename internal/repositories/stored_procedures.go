package repositories

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
	"github.com/pikami/cosmium/internal/resourceid"
	"golang.org/x/exp/maps"
)

func (r *DataRepository) GetAllStoredProcedures(databaseId string, collectionId string) ([]repositorymodels.StoredProcedure, repositorymodels.RepositoryStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	return maps.Values(r.storeState.StoredProcedures[databaseId][collectionId]), repositorymodels.StatusOk
}

func (r *DataRepository) GetStoredProcedure(databaseId string, collectionId string, spId string) (repositorymodels.StoredProcedure, repositorymodels.RepositoryStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return repositorymodels.StoredProcedure{}, repositorymodels.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return repositorymodels.StoredProcedure{}, repositorymodels.StatusNotFound
	}

	if sp, ok := r.storeState.StoredProcedures[databaseId][collectionId][spId]; ok {
		return sp, repositorymodels.StatusOk
	}

	return repositorymodels.StoredProcedure{}, repositorymodels.StatusNotFound
}

func (r *DataRepository) DeleteStoredProcedure(databaseId string, collectionId string, spId string) repositorymodels.RepositoryStatus {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return repositorymodels.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return repositorymodels.StatusNotFound
	}

	if _, ok := r.storeState.StoredProcedures[databaseId][collectionId][spId]; !ok {
		return repositorymodels.StatusNotFound
	}

	delete(r.storeState.StoredProcedures[databaseId][collectionId], spId)

	return repositorymodels.StatusOk
}

func (r *DataRepository) CreateStoredProcedure(databaseId string, collectionId string, sp repositorymodels.StoredProcedure) (repositorymodels.StoredProcedure, repositorymodels.RepositoryStatus) {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	var ok bool
	var database repositorymodels.Database
	var collection repositorymodels.Collection
	if sp.ID == "" {
		return repositorymodels.StoredProcedure{}, repositorymodels.BadRequest
	}

	if database, ok = r.storeState.Databases[databaseId]; !ok {
		return repositorymodels.StoredProcedure{}, repositorymodels.StatusNotFound
	}

	if collection, ok = r.storeState.Collections[databaseId][collectionId]; !ok {
		return repositorymodels.StoredProcedure{}, repositorymodels.StatusNotFound
	}

	if _, ok = r.storeState.StoredProcedures[databaseId][collectionId][sp.ID]; ok {
		return repositorymodels.StoredProcedure{}, repositorymodels.Conflict
	}

	sp.TimeStamp = time.Now().Unix()
	sp.ResourceID = resourceid.NewCombined(database.ResourceID, collection.ResourceID, resourceid.New())
	sp.ETag = fmt.Sprintf("\"%s\"", uuid.New())
	sp.Self = fmt.Sprintf("dbs/%s/colls/%s/sprocs/%s/", database.ResourceID, collection.ResourceID, sp.ResourceID)

	r.storeState.StoredProcedures[databaseId][collectionId][sp.ID] = sp

	return sp, repositorymodels.StatusOk
}
