package datastore

type DataStore interface {
	GetAllDatabases() ([]Database, DataStoreStatus)
	GetDatabase(databaseId string) (Database, DataStoreStatus)
	DeleteDatabase(databaseId string) DataStoreStatus
	CreateDatabase(newDatabase Database) (Database, DataStoreStatus)

	GetAllCollections(databaseId string) ([]Collection, DataStoreStatus)
	GetCollection(databaseId string, collectionId string) (Collection, DataStoreStatus)
	DeleteCollection(databaseId string, collectionId string) DataStoreStatus
	CreateCollection(databaseId string, newCollection Collection) (Collection, DataStoreStatus)

	GetAllDocuments(databaseId string, collectionId string) ([]Document, DataStoreStatus)
	GetDocumentIterator(databaseId string, collectionId string) (DocumentIterator, DataStoreStatus)
	GetDocument(databaseId string, collectionId string, documentId string) (Document, DataStoreStatus)
	DeleteDocument(databaseId string, collectionId string, documentId string) DataStoreStatus
	CreateDocument(databaseId string, collectionId string, document map[string]interface{}) (Document, DataStoreStatus)

	GetAllTriggers(databaseId string, collectionId string) ([]Trigger, DataStoreStatus)
	GetTrigger(databaseId string, collectionId string, triggerId string) (Trigger, DataStoreStatus)
	DeleteTrigger(databaseId string, collectionId string, triggerId string) DataStoreStatus
	CreateTrigger(databaseId string, collectionId string, trigger Trigger) (Trigger, DataStoreStatus)

	GetAllStoredProcedures(databaseId string, collectionId string) ([]StoredProcedure, DataStoreStatus)
	GetStoredProcedure(databaseId string, collectionId string, storedProcedureId string) (StoredProcedure, DataStoreStatus)
	DeleteStoredProcedure(databaseId string, collectionId string, storedProcedureId string) DataStoreStatus
	CreateStoredProcedure(databaseId string, collectionId string, storedProcedure StoredProcedure) (StoredProcedure, DataStoreStatus)

	GetAllUserDefinedFunctions(databaseId string, collectionId string) ([]UserDefinedFunction, DataStoreStatus)
	GetUserDefinedFunction(databaseId string, collectionId string, udfId string) (UserDefinedFunction, DataStoreStatus)
	DeleteUserDefinedFunction(databaseId string, collectionId string, udfId string) DataStoreStatus
	CreateUserDefinedFunction(databaseId string, collectionId string, udf UserDefinedFunction) (UserDefinedFunction, DataStoreStatus)

	GetPartitionKeyRanges(databaseId string, collectionId string) ([]PartitionKeyRange, DataStoreStatus)

	Close()
	DumpToJson() (string, error)
}

type DocumentIterator interface {
	Next() (Document, DataStoreStatus)
	Close()
}
