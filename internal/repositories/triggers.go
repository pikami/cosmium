package repositories

import repositorymodels "github.com/pikami/cosmium/internal/repository_models"

func GetAllTriggers(databaseId string, collectionId string) ([]repositorymodels.Trigger, repositorymodels.RepositoryStatus) {
	return triggers, repositorymodels.StatusOk
}
