package repositorymodels

type Database struct {
	ID        string `json:"id"`
	TimeStamp int64  `json:"_ts"`
	UniqueID  string `json:"_rid"`
	ETag      string `json:"_etag"`
}

type RepositoryStatus int

const (
	StatusOk       = 1
	StatusNotFound = 2
	Conflict       = 3
	BadRequest     = 4
)

type Collection struct {
	ID             string                   `json:"id"`
	IndexingPolicy CollectionIndexingPolicy `json:"indexingPolicy"`
	PartitionKey   CollectionPartitionKey   `json:"partitionKey"`
	UniqueID       string                   `json:"_rid"`
	TimeStamp      int64                    `json:"_ts"`
	Self           string                   `json:"_self"`
	ETag           string                   `json:"_etag"`
	Docs           string                   `json:"_docs"`
	Sprocs         string                   `json:"_sprocs"`
	Triggers       string                   `json:"_triggers"`
	Udfs           string                   `json:"_udfs"`
	Conflicts      string                   `json:"_conflicts"`
}

type CollectionIndexingPolicy struct {
	IndexingMode  string                         `json:"indexingMode"`
	Automatic     bool                           `json:"automatic"`
	IncludedPaths []CollectionIndexingPolicyPath `json:"includedPaths"`
	ExcludedPaths []CollectionIndexingPolicyPath `json:"excludedPaths"`
}

type CollectionIndexingPolicyPath struct {
	Path    string `json:"path"`
	Indexes []struct {
		Kind      string `json:"kind"`
		DataType  string `json:"dataType"`
		Precision int    `json:"precision"`
	} `json:"indexes"`
}

type CollectionPartitionKey struct {
	Paths   []string `json:"paths"`
	Kind    string   `json:"kind"`
	Version int      `json:"Version"`
}

type UserDefinedFunction struct {
	Body string `json:"body"`
	ID   string `json:"id"`
	Rid  string `json:"_rid"`
	Ts   int    `json:"_ts"`
	Self string `json:"_self"`
	Etag string `json:"_etag"`
}

type StoredProcedure struct {
	Body string `json:"body"`
	ID   string `json:"id"`
	Rid  string `json:"_rid"`
	Ts   int    `json:"_ts"`
	Self string `json:"_self"`
	Etag string `json:"_etag"`
}

type Trigger struct {
	Body             string `json:"body"`
	ID               string `json:"id"`
	TriggerOperation string `json:"triggerOperation"`
	TriggerType      string `json:"triggerType"`
	Rid              string `json:"_rid"`
	Ts               int    `json:"_ts"`
	Self             string `json:"_self"`
	Etag             string `json:"_etag"`
}

type Document map[string]interface{}

type PartitionKeyRange struct {
	Rid                string `json:"_rid"`
	ID                 string `json:"id"`
	Etag               string `json:"_etag"`
	MinInclusive       string `json:"minInclusive"`
	MaxExclusive       string `json:"maxExclusive"`
	RidPrefix          int    `json:"ridPrefix"`
	Self               string `json:"_self"`
	ThroughputFraction int    `json:"throughputFraction"`
	Status             string `json:"status"`
	Parents            []any  `json:"parents"`
	Ts                 int    `json:"_ts"`
	Lsn                int    `json:"lsn"`
}

type State struct {
	// Map databaseId -> Database
	Databases map[string]Database `json:"databases"`

	// Map databaseId -> collectionId -> Collection
	Collections map[string]map[string]Collection `json:"collections"`

	// Map databaseId -> collectionId -> documentId -> Documents
	Documents map[string]map[string]map[string]Document `json:"documents"`
}
