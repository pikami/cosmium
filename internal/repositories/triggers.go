package repositories

import repositorymodels "github.com/pikami/cosmium/internal/repository_models"

var triggers = []repositorymodels.Trigger{}

func GetAllTriggers(databaseId string, collectionId string) ([]repositorymodels.Trigger, repositorymodels.RepositoryStatus) {
	filteredTriggers := make([]repositorymodels.Trigger, 0)

	for _, coll := range triggers {
		if coll.Internals.DatabaseId == databaseId && coll.Internals.CollectionId == collectionId {
			filteredTriggers = append(filteredTriggers, coll)
		}
	}

	return filteredTriggers, repositorymodels.StatusOk
}
