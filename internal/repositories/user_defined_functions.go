package repositories

var userDefinedFunctions = []UserDefinedFunction{}

func GetAllUserDefinedFunctions(databaseId string, collectionId string) ([]UserDefinedFunction, RepositoryStatus) {
	udfs := make([]UserDefinedFunction, 0)

	for _, coll := range userDefinedFunctions {
		if coll.internals.databaseId == databaseId && coll.internals.collectionId == collectionId {
			udfs = append(udfs, coll)
		}
	}

	return udfs, StatusOk
}
