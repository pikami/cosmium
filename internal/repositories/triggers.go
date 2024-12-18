package repositories

import repositorymodels "github.com/pikami/cosmium/internal/repository_models"

func (r *DataRepository) GetAllTriggers(databaseId string, collectionId string) ([]repositorymodels.Trigger, repositorymodels.RepositoryStatus) {
	return r.triggers, repositorymodels.StatusOk
}
