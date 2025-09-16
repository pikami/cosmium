package datastore

type InitialDataModel struct {
	// Map databaseId -> Database
	Databases map[string]Database `json:"databases"`

	// Map databaseId -> collectionId -> Collection
	Collections map[string]map[string]Collection `json:"collections"`

	// Map databaseId -> collectionId -> documentId -> Documents
	Documents map[string]map[string]map[string]Document `json:"documents"`

	// Map databaseId -> collectionId -> triggerId -> Trigger
	Triggers map[string]map[string]map[string]Trigger `json:"triggers"`

	// Map databaseId -> collectionId -> spId -> StoredProcedure
	StoredProcedures map[string]map[string]map[string]StoredProcedure `json:"sprocs"`

	// Map databaseId -> collectionId -> udfId -> UserDefinedFunction
	UserDefinedFunctions map[string]map[string]map[string]UserDefinedFunction `json:"udfs"`
}
