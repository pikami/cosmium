package repositories

import repositorymodels "github.com/pikami/cosmium/internal/repository_models"

var userDefinedFunctions = []repositorymodels.UserDefinedFunction{}

func GetAllUserDefinedFunctions(databaseId string, collectionId string) ([]repositorymodels.UserDefinedFunction, repositorymodels.RepositoryStatus) {
	udfs := make([]repositorymodels.UserDefinedFunction, 0)

	for _, coll := range userDefinedFunctions {
		if coll.Internals.DatabaseId == databaseId && coll.Internals.CollectionId == collectionId {
			udfs = append(udfs, coll)
		}
	}

	return udfs, repositorymodels.StatusOk
}
