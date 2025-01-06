package repositories

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
	"github.com/pikami/cosmium/internal/resourceid"
	"golang.org/x/exp/maps"
)

func (r *DataRepository) GetAllDatabases() ([]repositorymodels.Database, repositorymodels.RepositoryStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	return maps.Values(r.storeState.Databases), repositorymodels.StatusOk
}

func (r *DataRepository) GetDatabase(id string) (repositorymodels.Database, repositorymodels.RepositoryStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	if database, ok := r.storeState.Databases[id]; ok {
		return database, repositorymodels.StatusOk
	}

	return repositorymodels.Database{}, repositorymodels.StatusNotFound
}

func (r *DataRepository) DeleteDatabase(id string) repositorymodels.RepositoryStatus {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	if _, ok := r.storeState.Databases[id]; !ok {
		return repositorymodels.StatusNotFound
	}

	delete(r.storeState.Databases, id)

	return repositorymodels.StatusOk
}

func (r *DataRepository) CreateDatabase(newDatabase repositorymodels.Database) (repositorymodels.Database, repositorymodels.RepositoryStatus) {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	if _, ok := r.storeState.Databases[newDatabase.ID]; ok {
		return repositorymodels.Database{}, repositorymodels.Conflict
	}

	newDatabase.TimeStamp = time.Now().Unix()
	newDatabase.ResourceID = resourceid.New()
	newDatabase.ETag = fmt.Sprintf("\"%s\"", uuid.New())
	newDatabase.Self = fmt.Sprintf("dbs/%s/", newDatabase.ResourceID)

	r.storeState.Databases[newDatabase.ID] = newDatabase
	r.storeState.Collections[newDatabase.ID] = make(map[string]repositorymodels.Collection)
	r.storeState.Documents[newDatabase.ID] = make(map[string]map[string]repositorymodels.Document)
	r.storeState.Triggers[newDatabase.ID] = make(map[string]map[string]repositorymodels.Trigger)
	r.storeState.StoredProcedures[newDatabase.ID] = make(map[string]map[string]repositorymodels.StoredProcedure)
	r.storeState.UserDefinedFunctions[newDatabase.ID] = make(map[string]map[string]repositorymodels.UserDefinedFunction)

	return newDatabase, repositorymodels.StatusOk
}
