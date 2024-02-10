package repositories

var storedProcedures = []StoredProcedure{}

func GetAllStoredProcedures(databaseId string, collectionId string) ([]StoredProcedure, RepositoryStatus) {
	sps := make([]StoredProcedure, 0)

	for _, coll := range storedProcedures {
		if coll.internals.databaseId == databaseId && coll.internals.collectionId == collectionId {
			sps = append(sps, coll)
		}
	}

	return sps, StatusOk
}
