package repositories

func GetPartitionKeyRanges(databaseId string, collectionId string) ([]PartitionKeyRange, RepositoryStatus) {
	// I have no idea what this is tbh
	return []PartitionKeyRange{
		{
			Rid:                "ZxlyAP7rKwACAAAAAAAAUA==",
			ID:                 "0",
			Etag:               "\"00005504-0000-0100-0000-65c555490000\"",
			MinInclusive:       "",
			MaxExclusive:       "FF",
			RidPrefix:          0,
			Self:               "dbs/ZxlyAA==/colls/ZxlyAP7rKwA=/pkranges/ZxlyAP7rKwACAAAAAAAAUA==/",
			ThroughputFraction: 1,
			Status:             "online",
			Parents:            []interface{}{},
			Ts:                 1707431241,
			Lsn:                17,
		},
	}, StatusOk
}
