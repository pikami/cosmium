package repositories

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
	"github.com/pikami/cosmium/internal/resourceid"
	"golang.org/x/exp/maps"
)

func GetAllDatabases() ([]repositorymodels.Database, repositorymodels.RepositoryStatus) {
	return maps.Values(storeState.Databases), repositorymodels.StatusOk
}

func GetDatabase(id string) (repositorymodels.Database, repositorymodels.RepositoryStatus) {
	if database, ok := storeState.Databases[id]; ok {
		return database, repositorymodels.StatusOk
	}

	return repositorymodels.Database{}, repositorymodels.StatusNotFound
}

func DeleteDatabase(id string) repositorymodels.RepositoryStatus {
	if _, ok := storeState.Databases[id]; !ok {
		return repositorymodels.StatusNotFound
	}

	delete(storeState.Databases, id)

	return repositorymodels.StatusOk
}

func CreateDatabase(newDatabase repositorymodels.Database) (repositorymodels.Database, repositorymodels.RepositoryStatus) {
	if _, ok := storeState.Databases[newDatabase.ID]; ok {
		return repositorymodels.Database{}, repositorymodels.Conflict
	}

	newDatabase.TimeStamp = time.Now().Unix()
	newDatabase.ResourceID = resourceid.New()
	newDatabase.ETag = fmt.Sprintf("\"%s\"", uuid.New())
	newDatabase.Self = fmt.Sprintf("dbs/%s/", newDatabase.ResourceID)

	storeState.Databases[newDatabase.ID] = newDatabase
	storeState.Collections[newDatabase.ID] = make(map[string]repositorymodels.Collection)
	storeState.Documents[newDatabase.ID] = make(map[string]map[string]repositorymodels.Document)

	return newDatabase, repositorymodels.StatusOk
}
