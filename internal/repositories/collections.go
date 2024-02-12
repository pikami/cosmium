package repositories

import (
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
	structhidrators "github.com/pikami/cosmium/internal/struct_hidrators"
)

var collections = []repositorymodels.Collection{
	{ID: "db1"},
	{ID: "db2"},
}

func GetAllCollections(databaseId string) ([]repositorymodels.Collection, repositorymodels.RepositoryStatus) {
	dbCollections := make([]repositorymodels.Collection, 0)

	for _, coll := range collections {
		if coll.Internals.DatabaseId == databaseId {
			dbCollections = append(dbCollections, coll)
		}
	}

	return dbCollections, repositorymodels.StatusOk
}

func GetCollection(databaseId string, id string) (repositorymodels.Collection, repositorymodels.RepositoryStatus) {
	for _, coll := range collections {
		if coll.Internals.DatabaseId == databaseId && coll.ID == id {
			return coll, repositorymodels.StatusOk
		}
	}

	return repositorymodels.Collection{}, repositorymodels.StatusNotFound
}

func DeleteCollection(databaseId string, id string) repositorymodels.RepositoryStatus {
	for index, coll := range collections {
		if coll.Internals.DatabaseId == databaseId && coll.ID == id {
			collections = append(collections[:index], collections[index+1:]...)
			return repositorymodels.StatusOk
		}
	}

	return repositorymodels.StatusNotFound
}

func CreateCollection(databaseId string, newCollection repositorymodels.Collection) (repositorymodels.Collection, repositorymodels.RepositoryStatus) {
	for _, coll := range collections {
		if coll.Internals.DatabaseId == databaseId && coll.ID == newCollection.ID {
			return repositorymodels.Collection{}, repositorymodels.Conflict
		}
	}

	newCollection = structhidrators.Hidrate(newCollection).(repositorymodels.Collection)

	newCollection.Internals = struct{ DatabaseId string }{
		DatabaseId: databaseId,
	}
	collections = append(collections, newCollection)
	return newCollection, repositorymodels.StatusOk
}
