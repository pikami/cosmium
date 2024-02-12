package repositories

import repositorymodels "github.com/pikami/cosmium/internal/repository_models"

var storedProcedures = []repositorymodels.StoredProcedure{}

func GetAllStoredProcedures(databaseId string, collectionId string) ([]repositorymodels.StoredProcedure, repositorymodels.RepositoryStatus) {
	sps := make([]repositorymodels.StoredProcedure, 0)

	for _, coll := range storedProcedures {
		if coll.Internals.DatabaseId == databaseId && coll.Internals.CollectionId == collectionId {
			sps = append(sps, coll)
		}
	}

	return sps, repositorymodels.StatusOk
}
