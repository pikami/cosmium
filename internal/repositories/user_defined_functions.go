package repositories

import repositorymodels "github.com/pikami/cosmium/internal/repository_models"

func GetAllUserDefinedFunctions(databaseId string, collectionId string) ([]repositorymodels.UserDefinedFunction, repositorymodels.RepositoryStatus) {
	return userDefinedFunctions, repositorymodels.StatusOk
}
