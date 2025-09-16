package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	jsonpatch "github.com/cosmiumdev/json-patch/v5"
	"github.com/gin-gonic/gin"
	apimodels "github.com/pikami/cosmium/api/api_models"
	"github.com/pikami/cosmium/api/headers"
	"github.com/pikami/cosmium/internal/constants"
	"github.com/pikami/cosmium/internal/converters"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/logger"
	"github.com/pikami/cosmium/parsers"
	"github.com/pikami/cosmium/parsers/nosql"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
)

func (h *Handlers) GetAllDocuments(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	documents, status := h.dataStore.GetAllDocuments(databaseId, collectionId)
	if status == datastore.StatusOk {
		collection, _ := h.dataStore.GetCollection(databaseId, collectionId)

		c.Header(headers.ItemCount, fmt.Sprintf("%d", len(documents)))
		c.IndentedJSON(http.StatusOK, gin.H{
			"_rid":      collection.ID,
			"Documents": documents,
			"_count":    len(documents),
		})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}

func (h *Handlers) GetDocument(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	documentId := c.Param("docId")

	document, status := h.dataStore.GetDocument(databaseId, collectionId, documentId)
	if status == datastore.StatusOk {
		c.IndentedJSON(http.StatusOK, document)
		return
	}

	if status == datastore.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, constants.NotFoundResponse)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}

func (h *Handlers) DeleteDocument(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	documentId := c.Param("docId")

	status := h.dataStore.DeleteDocument(databaseId, collectionId, documentId)
	if status == datastore.StatusOk {
		c.Status(http.StatusNoContent)
		return
	}

	if status == datastore.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, constants.NotFoundResponse)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}

// TODO: Maybe move "replace" logic to data store
func (h *Handlers) ReplaceDocument(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	documentId := c.Param("docId")

	var requestBody map[string]interface{}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	status := h.dataStore.DeleteDocument(databaseId, collectionId, documentId)
	if status == datastore.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, constants.NotFoundResponse)
		return
	}

	createdDocument, status := h.dataStore.CreateDocument(databaseId, collectionId, requestBody)
	if status == datastore.Conflict {
		c.IndentedJSON(http.StatusConflict, constants.ConflictResponse)
		return
	}

	if status == datastore.StatusOk {
		c.IndentedJSON(http.StatusCreated, createdDocument)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}

func (h *Handlers) PatchDocument(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	documentId := c.Param("docId")

	document, status := h.dataStore.GetDocument(databaseId, collectionId, documentId)
	if status == datastore.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, constants.NotFoundResponse)
		return
	}

	var requestBody map[string]interface{}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	operations := requestBody["operations"]
	operationsBytes, err := json.Marshal(operations)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not decode operations"})
		return
	}

	patch, err := jsonpatch.DecodePatch(operationsBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	currentDocumentBytes, err := json.Marshal(document)
	if err != nil {
		logger.ErrorLn("Failed to marshal existing document:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to marshal existing document"})
		return
	}

	modifiedDocumentBytes, err := patch.Apply(currentDocumentBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var modifiedDocument map[string]interface{}
	err = json.Unmarshal(modifiedDocumentBytes, &modifiedDocument)
	if err != nil {
		logger.ErrorLn("Failed to unmarshal modified document:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to unmarshal modified document"})
		return
	}

	if modifiedDocument["id"] != document["id"] {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "The ID field cannot be modified"})
		return
	}

	status = h.dataStore.DeleteDocument(databaseId, collectionId, documentId)
	if status == datastore.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, constants.NotFoundResponse)
		return
	}

	createdDocument, status := h.dataStore.CreateDocument(databaseId, collectionId, modifiedDocument)
	if status == datastore.Conflict {
		c.IndentedJSON(http.StatusConflict, constants.ConflictResponse)
		return
	}

	if status == datastore.StatusOk {
		c.IndentedJSON(http.StatusCreated, createdDocument)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}

func (h *Handlers) DocumentsPost(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	// Handle batch requests
	isBatchRequest, _ := strconv.ParseBool(c.GetHeader(headers.IsBatchRequest))
	if isBatchRequest {
		h.handleBatchRequest(c)
		return
	}

	var requestBody map[string]interface{}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Handle query plan requests
	isQueryPlanRequest, _ := strconv.ParseBool(c.GetHeader(headers.IsQueryPlanRequest))
	if isQueryPlanRequest {
		c.IndentedJSON(http.StatusOK, constants.QueryPlanResponse)
		return
	}

	// Handle query requests
	isQueryRequest, _ := strconv.ParseBool(c.GetHeader(headers.IsQuery))
	isQueryRequestAltHeader, _ := strconv.ParseBool(c.GetHeader(headers.Query))
	if isQueryRequest || isQueryRequestAltHeader {
		h.handleDocumentQuery(c, requestBody)
		return
	}

	if requestBody["id"] == "" {
		c.JSON(http.StatusBadRequest, constants.BadRequestResponse)
		return
	}

	isUpsert, _ := strconv.ParseBool(c.GetHeader(headers.IsUpsert))
	if isUpsert {
		h.dataStore.DeleteDocument(databaseId, collectionId, requestBody["id"].(string))
	}

	createdDocument, status := h.dataStore.CreateDocument(databaseId, collectionId, requestBody)
	if status == datastore.Conflict {
		c.IndentedJSON(http.StatusConflict, constants.ConflictResponse)
		return
	}

	if status == datastore.StatusOk {
		c.IndentedJSON(http.StatusCreated, createdDocument)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}

func parametersToMap(pairs []interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for _, pair := range pairs {
		if pairMap, ok := pair.(map[string]interface{}); ok {
			result[pairMap["name"].(string)] = pairMap["value"]
		}
	}

	return result
}

func (h *Handlers) handleDocumentQuery(c *gin.Context, requestBody map[string]interface{}) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	var queryParameters map[string]interface{}
	if paramsArray, ok := requestBody["parameters"].([]interface{}); ok {
		queryParameters = parametersToMap(paramsArray)
	}

	queryText := requestBody["query"].(string)
	docs, status := h.executeQueryDocuments(databaseId, collectionId, queryText, queryParameters)
	if status != datastore.StatusOk {
		// TODO: Currently we return everything if the query fails
		h.GetAllDocuments(c)
		return
	}

	collection, _ := h.dataStore.GetCollection(databaseId, collectionId)
	c.Header(headers.ItemCount, fmt.Sprintf("%d", len(docs)))
	c.IndentedJSON(http.StatusOK, gin.H{
		"_rid":      collection.ResourceID,
		"Documents": docs,
		"_count":    len(docs),
	})
}

func (h *Handlers) handleBatchRequest(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	batchOperations := make([]apimodels.BatchOperation, 0)
	if err := c.BindJSON(&batchOperations); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	batchOperationResults := make([]apimodels.BatchOperationResult, len(batchOperations))
	for idx, operation := range batchOperations {
		switch operation.OperationType {
		case apimodels.BatchOperationTypeCreate:
			createdDocument, status := h.dataStore.CreateDocument(databaseId, collectionId, operation.ResourceBody)
			responseCode := dataStoreStatusToResponseCode(status)
			if status == datastore.StatusOk {
				responseCode = http.StatusCreated
			}
			batchOperationResults[idx] = apimodels.BatchOperationResult{
				StatusCode:   responseCode,
				ResourceBody: createdDocument,
			}
		case apimodels.BatchOperationTypeDelete:
			status := h.dataStore.DeleteDocument(databaseId, collectionId, operation.Id)
			responseCode := dataStoreStatusToResponseCode(status)
			if status == datastore.StatusOk {
				responseCode = http.StatusNoContent
			}
			batchOperationResults[idx] = apimodels.BatchOperationResult{
				StatusCode: responseCode,
			}
		case apimodels.BatchOperationTypeReplace:
			deleteStatus := h.dataStore.DeleteDocument(databaseId, collectionId, operation.Id)
			if deleteStatus == datastore.StatusNotFound {
				batchOperationResults[idx] = apimodels.BatchOperationResult{
					StatusCode: http.StatusNotFound,
				}
				continue
			}
			createdDocument, createStatus := h.dataStore.CreateDocument(databaseId, collectionId, operation.ResourceBody)
			responseCode := dataStoreStatusToResponseCode(createStatus)
			if createStatus == datastore.StatusOk {
				responseCode = http.StatusCreated
			}
			batchOperationResults[idx] = apimodels.BatchOperationResult{
				StatusCode:   responseCode,
				ResourceBody: createdDocument,
			}
		case apimodels.BatchOperationTypeUpsert:
			documentId := operation.ResourceBody["id"].(string)
			h.dataStore.DeleteDocument(databaseId, collectionId, documentId)
			createdDocument, createStatus := h.dataStore.CreateDocument(databaseId, collectionId, operation.ResourceBody)
			responseCode := dataStoreStatusToResponseCode(createStatus)
			if createStatus == datastore.StatusOk {
				responseCode = http.StatusCreated
			}
			batchOperationResults[idx] = apimodels.BatchOperationResult{
				StatusCode:   responseCode,
				ResourceBody: createdDocument,
			}
		case apimodels.BatchOperationTypeRead:
			document, status := h.dataStore.GetDocument(databaseId, collectionId, operation.Id)
			batchOperationResults[idx] = apimodels.BatchOperationResult{
				StatusCode:   dataStoreStatusToResponseCode(status),
				ResourceBody: document,
			}
		case apimodels.BatchOperationTypePatch:
			batchOperationResults[idx] = apimodels.BatchOperationResult{
				StatusCode: http.StatusNotImplemented,
				Message:    "Patch operation is not implemented",
			}
		default:
			batchOperationResults[idx] = apimodels.BatchOperationResult{
				StatusCode: http.StatusBadRequest,
				Message:    "Unknown operation type",
			}
		}
	}

	c.JSON(http.StatusOK, batchOperationResults)
}

func dataStoreStatusToResponseCode(status datastore.DataStoreStatus) int {
	switch status {
	case datastore.StatusOk:
		return http.StatusOK
	case datastore.StatusNotFound:
		return http.StatusNotFound
	case datastore.Conflict:
		return http.StatusConflict
	case datastore.BadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func (h *Handlers) executeQueryDocuments(databaseId string, collectionId string, query string, queryParameters map[string]interface{}) ([]memoryexecutor.RowType, datastore.DataStoreStatus) {
	parsedQuery, err := nosql.Parse("", []byte(query))
	if err != nil {
		logger.Errorf("Failed to parse query: %s\nerr: %v", query, err)
		return nil, datastore.BadRequest
	}

	allDocumentsIterator, status := h.dataStore.GetDocumentIterator(databaseId, collectionId)
	if status != datastore.StatusOk {
		return nil, status
	}
	defer allDocumentsIterator.Close()

	rowsIterator := converters.NewDocumentToRowTypeIterator(allDocumentsIterator)

	if typedQuery, ok := parsedQuery.(parsers.SelectStmt); ok {
		typedQuery.Parameters = queryParameters
		return memoryexecutor.ExecuteQuery(typedQuery, rowsIterator), datastore.StatusOk
	}

	return nil, datastore.BadRequest
}
