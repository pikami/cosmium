package repositories

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

var databases = []repositorymodels.Database{
	{ID: "db1"},
	{ID: "db2"},
}

func GetAllDatabases() ([]repositorymodels.Database, repositorymodels.RepositoryStatus) {
	return databases, repositorymodels.StatusOk
}

func GetDatabase(id string) (repositorymodels.Database, repositorymodels.RepositoryStatus) {
	for _, db := range databases {
		if db.ID == id {
			return db, repositorymodels.StatusOk
		}
	}

	return repositorymodels.Database{}, repositorymodels.StatusNotFound
}

func DeleteDatabase(id string) repositorymodels.RepositoryStatus {
	for index, db := range databases {
		if db.ID == id {
			databases = append(databases[:index], databases[index+1:]...)
			return repositorymodels.StatusOk
		}
	}

	return repositorymodels.StatusNotFound
}

func CreateDatabase(newDatabase repositorymodels.Database) repositorymodels.RepositoryStatus {
	for _, db := range databases {
		if db.ID == newDatabase.ID {
			return repositorymodels.Conflict
		}
	}

	newDatabase.TimeStamp = time.Now().Unix()
	newDatabase.UniqueID = uuid.New().String()
	newDatabase.ETag = fmt.Sprintf("\"%s\"", newDatabase.UniqueID)
	databases = append(databases, newDatabase)
	return repositorymodels.StatusOk
}
