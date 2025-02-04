package apimodels

const (
	BatchOperationTypeCreate  = "Create"
	BatchOperationTypeDelete  = "Delete"
	BatchOperationTypeReplace = "Replace"
	BatchOperationTypeUpsert  = "Upsert"
	BatchOperationTypeRead    = "Read"
	BatchOperationTypePatch   = "Patch"
)

type BatchOperation struct {
	OperationType string                 `json:"operationType"`
	Id            string                 `json:"id"`
	ResourceBody  map[string]interface{} `json:"resourceBody"`
}

type BatchOperationResult struct {
	StatusCode    int                    `json:"statusCode"`
	RequestCharge float64                `json:"requestCharge"`
	ResourceBody  map[string]interface{} `json:"resourceBody"`
	Etag          string                 `json:"etag"`
	Message       string                 `json:"message"`
}
