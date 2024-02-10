package repositories

var triggers = []Trigger{}

func GetAllTriggers(databaseId string, collectionId string) ([]Trigger, RepositoryStatus) {
	sps := make([]Trigger, 0)

	for _, coll := range triggers {
		if coll.internals.databaseId == databaseId && coll.internals.collectionId == collectionId {
			sps = append(sps, coll)
		}
	}

	return sps, StatusOk
}
