package repositories

import repositorymodels "github.com/pikami/cosmium/internal/repository_models"

func (r *DataRepository) GetAllUserDefinedFunctions(databaseId string, collectionId string) ([]repositorymodels.UserDefinedFunction, repositorymodels.RepositoryStatus) {
	return r.userDefinedFunctions, repositorymodels.StatusOk
}
