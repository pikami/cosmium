package datastore

type Database struct {
	ID         string `json:"id"`
	TimeStamp  int64  `json:"_ts"`
	ResourceID string `json:"_rid"`
	ETag       string `json:"_etag"`
	Self       string `json:"_self"`
}

type DataStoreStatus int

const (
	StatusOk       DataStoreStatus = 1
	StatusNotFound DataStoreStatus = 2
	Conflict       DataStoreStatus = 3
	BadRequest     DataStoreStatus = 4
	IterEOF        DataStoreStatus = 5
	Unknown        DataStoreStatus = 6
)

type TriggerOperation string

const (
	All     TriggerOperation = "All"
	Create  TriggerOperation = "Create"
	Delete  TriggerOperation = "Delete"
	Replace TriggerOperation = "Replace"
)

type TriggerType string

const (
	Pre  TriggerType = "Pre"
	Post TriggerType = "Post"
)

type Collection struct {
	ID             string                   `json:"id"`
	IndexingPolicy CollectionIndexingPolicy `json:"indexingPolicy"`
	PartitionKey   CollectionPartitionKey   `json:"partitionKey"`
	ResourceID     string                   `json:"_rid"`
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
	Body       string `json:"body"`
	ID         string `json:"id"`
	ResourceID string `json:"_rid"`
	TimeStamp  int64  `json:"_ts"`
	Self       string `json:"_self"`
	ETag       string `json:"_etag"`
}

type StoredProcedure struct {
	Body       string `json:"body"`
	ID         string `json:"id"`
	ResourceID string `json:"_rid"`
	TimeStamp  int64  `json:"_ts"`
	Self       string `json:"_self"`
	ETag       string `json:"_etag"`
}

type Trigger struct {
	Body             string           `json:"body"`
	ID               string           `json:"id"`
	TriggerOperation TriggerOperation `json:"triggerOperation"`
	TriggerType      TriggerType      `json:"triggerType"`
	ResourceID       string           `json:"_rid"`
	TimeStamp        int64            `json:"_ts"`
	Self             string           `json:"_self"`
	ETag             string           `json:"_etag"`
}

type Document map[string]interface{}

type PartitionKeyRange struct {
	ResourceID         string `json:"_rid"`
	ID                 string `json:"id"`
	Etag               string `json:"_etag"`
	MinInclusive       string `json:"minInclusive"`
	MaxExclusive       string `json:"maxExclusive"`
	RidPrefix          int    `json:"ridPrefix"`
	Self               string `json:"_self"`
	ThroughputFraction int    `json:"throughputFraction"`
	Status             string `json:"status"`
	Parents            []any  `json:"parents"`
	TimeStamp          int64  `json:"_ts"`
	Lsn                int    `json:"lsn"`
}
