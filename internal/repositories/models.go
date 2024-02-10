package repositories

type Database struct {
	ID string `json:"id"`
}

type RepositoryStatus int

const (
	StatusOk       = 1
	StatusNotFound = 2
	Conflict       = 3
	BadRequest     = 4
)

type Collection struct {
	ID             string `json:"id"`
	IndexingPolicy struct {
		IndexingMode  string `json:"indexingMode"`
		Automatic     bool   `json:"automatic"`
		IncludedPaths []struct {
			Path    string `json:"path"`
			Indexes []struct {
				Kind      string `json:"kind"`
				DataType  string `json:"dataType"`
				Precision int    `json:"precision"`
			} `json:"indexes"`
		} `json:"includedPaths"`
		ExcludedPaths []any `json:"excludedPaths"`
	} `json:"indexingPolicy"`
	PartitionKey struct {
		Paths   []string `json:"paths"`
		Kind    string   `json:"kind"`
		Version int      `json:"Version"`
	} `json:"partitionKey"`
	Rid       string `json:"_rid"`
	Ts        int    `json:"_ts"`
	Self      string `json:"_self"`
	Etag      string `json:"_etag"`
	Docs      string `json:"_docs"`
	Sprocs    string `json:"_sprocs"`
	Triggers  string `json:"_triggers"`
	Udfs      string `json:"_udfs"`
	Conflicts string `json:"_conflicts"`
	internals struct {
		databaseId string
	}
}

type UserDefinedFunction struct {
	Body      string `json:"body"`
	ID        string `json:"id"`
	Rid       string `json:"_rid"`
	Ts        int    `json:"_ts"`
	Self      string `json:"_self"`
	Etag      string `json:"_etag"`
	internals struct {
		databaseId   string
		collectionId string
	}
}

type StoredProcedure struct {
	Body      string `json:"body"`
	ID        string `json:"id"`
	Rid       string `json:"_rid"`
	Ts        int    `json:"_ts"`
	Self      string `json:"_self"`
	Etag      string `json:"_etag"`
	internals struct {
		databaseId   string
		collectionId string
	}
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
	internals        struct {
		databaseId   string
		collectionId string
	}
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
