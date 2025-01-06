package repositories

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
	"github.com/pikami/cosmium/internal/resourceid"
	"golang.org/x/exp/maps"
)

func (r *DataRepository) GetAllUserDefinedFunctions(databaseId string, collectionId string) ([]repositorymodels.UserDefinedFunction, repositorymodels.RepositoryStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	return maps.Values(r.storeState.UserDefinedFunctions[databaseId][collectionId]), repositorymodels.StatusOk
}

func (r *DataRepository) GetUserDefinedFunction(databaseId string, collectionId string, udfId string) (repositorymodels.UserDefinedFunction, repositorymodels.RepositoryStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return repositorymodels.UserDefinedFunction{}, repositorymodels.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return repositorymodels.UserDefinedFunction{}, repositorymodels.StatusNotFound
	}

	if udf, ok := r.storeState.UserDefinedFunctions[databaseId][collectionId][udfId]; ok {
		return udf, repositorymodels.StatusOk
	}

	return repositorymodels.UserDefinedFunction{}, repositorymodels.StatusNotFound
}

func (r *DataRepository) DeleteUserDefinedFunction(databaseId string, collectionId string, udfId string) repositorymodels.RepositoryStatus {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return repositorymodels.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return repositorymodels.StatusNotFound
	}

	if _, ok := r.storeState.UserDefinedFunctions[databaseId][collectionId][udfId]; !ok {
		return repositorymodels.StatusNotFound
	}

	delete(r.storeState.UserDefinedFunctions[databaseId][collectionId], udfId)

	return repositorymodels.StatusOk
}

func (r *DataRepository) CreateUserDefinedFunction(databaseId string, collectionId string, udf repositorymodels.UserDefinedFunction) (repositorymodels.UserDefinedFunction, repositorymodels.RepositoryStatus) {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	var ok bool
	var database repositorymodels.Database
	var collection repositorymodels.Collection
	if udf.ID == "" {
		return repositorymodels.UserDefinedFunction{}, repositorymodels.BadRequest
	}

	if database, ok = r.storeState.Databases[databaseId]; !ok {
		return repositorymodels.UserDefinedFunction{}, repositorymodels.StatusNotFound
	}

	if collection, ok = r.storeState.Collections[databaseId][collectionId]; !ok {
		return repositorymodels.UserDefinedFunction{}, repositorymodels.StatusNotFound
	}

	if _, ok := r.storeState.UserDefinedFunctions[databaseId][collectionId][udf.ID]; ok {
		return repositorymodels.UserDefinedFunction{}, repositorymodels.Conflict
	}

	udf.TimeStamp = time.Now().Unix()
	udf.ResourceID = resourceid.NewCombined(database.ResourceID, collection.ResourceID, resourceid.New())
	udf.ETag = fmt.Sprintf("\"%s\"", uuid.New())
	udf.Self = fmt.Sprintf("dbs/%s/colls/%s/udfs/%s/", database.ResourceID, collection.ResourceID, udf.ResourceID)

	r.storeState.UserDefinedFunctions[databaseId][collectionId][udf.ID] = udf

	return udf, repositorymodels.StatusOk
}
