package repositories

var collections = []Collection{
	{ID: "db1"},
	{ID: "db2"},
}

func GetAllCollections(databaseId string) ([]Collection, RepositoryStatus) {
	var dbCollections []Collection

	for _, coll := range collections {
		if coll.internals.databaseId == databaseId {
			dbCollections = append(dbCollections, coll)
		}
	}

	return dbCollections, StatusOk
}

func GetCollection(databaseId string, id string) (Collection, RepositoryStatus) {
	for _, coll := range collections {
		if coll.internals.databaseId == databaseId && coll.ID == id {
			return coll, StatusOk
		}
	}

	return Collection{}, StatusNotFound
}

func DeleteCollection(databaseId string, id string) RepositoryStatus {
	for index, coll := range collections {
		if coll.internals.databaseId == databaseId && coll.ID == id {
			collections = append(collections[:index], collections[index+1:]...)
			return StatusOk
		}
	}

	return StatusNotFound
}

func CreateCollection(databaseId string, newCollection Collection) RepositoryStatus {
	for _, coll := range collections {
		if coll.internals.databaseId == databaseId && coll.ID == newCollection.ID {
			return Conflict
		}
	}

	newCollection.internals = struct{ databaseId string }{
		databaseId: databaseId,
	}
	collections = append(collections, newCollection)
	return StatusOk
}
