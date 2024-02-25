package repositories

import repositorymodels "github.com/pikami/cosmium/internal/repository_models"

func GetAllStoredProcedures(databaseId string, collectionId string) ([]repositorymodels.StoredProcedure, repositorymodels.RepositoryStatus) {
	return storedProcedures, repositorymodels.StatusOk
}
