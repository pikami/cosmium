package repositories

var triggers = []Trigger{}

func GetAllTriggers(databaseId string, collectionId string) ([]Trigger, RepositoryStatus) {
	filteredTriggers := make([]Trigger, 0)

	for _, coll := range triggers {
		if coll.internals.databaseId == databaseId && coll.internals.collectionId == collectionId {
			filteredTriggers = append(filteredTriggers, coll)
		}
	}

	return filteredTriggers, StatusOk
}
