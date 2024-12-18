package repositories

import repositorymodels "github.com/pikami/cosmium/internal/repository_models"

func (r *DataRepository) GetAllStoredProcedures(databaseId string, collectionId string) ([]repositorymodels.StoredProcedure, repositorymodels.RepositoryStatus) {
	return r.storedProcedures, repositorymodels.StatusOk
}
