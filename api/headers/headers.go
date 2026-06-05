package headers

const (
	AIM                = "A-Im"
	Authorization      = "authorization"
	CosmosLsn          = "x-ms-cosmos-llsn"
	ETag               = "etag"
	GlobalCommittedLsn = "x-ms-global-committed-lsn"
	IfMatch            = "if-match"
	IfNoneMatch        = "if-none-match"
	IsBatchRequest     = "x-ms-cosmos-is-batch-request"
	IsQueryPlanRequest = "x-ms-cosmos-is-query-plan-request"
	IsUpsert           = "x-ms-documentdb-is-upsert"
	ItemCount          = "x-ms-item-count"
	LSN                = "lsn"
	XDate              = "x-ms-date"
	MaxItemCount       = "x-ms-max-item-count"
	ContinuationToken  = "x-ms-continuation"

	// Kinda retarded, but what can I do ¯\_(ツ)_/¯
	IsQuery = "x-ms-documentdb-isquery" // Sent from python sdk and web explorer
	Query   = "x-ms-documentdb-query"   // Sent from Go sdk
)
