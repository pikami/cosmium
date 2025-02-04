package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	jsonpatch "github.com/cosmiumdev/json-patch/v5"
	"github.com/gin-gonic/gin"
	apimodels "github.com/pikami/cosmium/api/api_models"
	"github.com/pikami/cosmium/internal/constants"
	"github.com/pikami/cosmium/internal/logger"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

func (h *Handlers) GetAllDocuments(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	documents, status := h.repository.GetAllDocuments(databaseId, collectionId)
	if status == repositorymodels.StatusOk {
		collection, _ := h.repository.GetCollection(databaseId, collectionId)

		c.Header("x-ms-item-count", fmt.Sprintf("%d", len(documents)))
		c.IndentedJSON(http.StatusOK, gin.H{
			"_rid":      collection.ID,
			"Documents": documents,
			"_count":    len(documents),
		})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) GetDocument(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	documentId := c.Param("docId")

	document, status := h.repository.GetDocument(databaseId, collectionId, documentId)
	if status == repositorymodels.StatusOk {
		c.IndentedJSON(http.StatusOK, document)
		return
	}

	if status == repositorymodels.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) DeleteDocument(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	documentId := c.Param("docId")

	status := h.repository.DeleteDocument(databaseId, collectionId, documentId)
	if status == repositorymodels.StatusOk {
		c.Status(http.StatusNoContent)
		return
	}

	if status == repositorymodels.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

// TODO: Maybe move "replace" logic to repository
func (h *Handlers) ReplaceDocument(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	documentId := c.Param("docId")

	var requestBody map[string]interface{}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	status := h.repository.DeleteDocument(databaseId, collectionId, documentId)
	if status == repositorymodels.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	createdDocument, status := h.repository.CreateDocument(databaseId, collectionId, requestBody)
	if status == repositorymodels.Conflict {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Conflict"})
		return
	}

	if status == repositorymodels.StatusOk {
		c.IndentedJSON(http.StatusCreated, createdDocument)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) PatchDocument(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	documentId := c.Param("docId")

	document, status := h.repository.GetDocument(databaseId, collectionId, documentId)
	if status == repositorymodels.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
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

	status = h.repository.DeleteDocument(databaseId, collectionId, documentId)
	if status == repositorymodels.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	createdDocument, status := h.repository.CreateDocument(databaseId, collectionId, modifiedDocument)
	if status == repositorymodels.Conflict {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Conflict"})
		return
	}

	if status == repositorymodels.StatusOk {
		c.IndentedJSON(http.StatusCreated, createdDocument)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) DocumentsPost(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	// Handle batch requests
	isBatchRequest, _ := strconv.ParseBool(c.GetHeader("x-ms-cosmos-is-batch-request"))
	if isBatchRequest {
		h.handleBatchRequest(c)
		return
	}

	var requestBody map[string]interface{}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	query := requestBody["query"]
	if query != nil {
		h.handleDocumentQuery(c, requestBody)
		return
	}

	if requestBody["id"] == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	isUpsert, _ := strconv.ParseBool(c.GetHeader("x-ms-documentdb-is-upsert"))
	if isUpsert {
		h.repository.DeleteDocument(databaseId, collectionId, requestBody["id"].(string))
	}

	createdDocument, status := h.repository.CreateDocument(databaseId, collectionId, requestBody)
	if status == repositorymodels.Conflict {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Conflict"})
		return
	}

	if status == repositorymodels.StatusOk {
		c.IndentedJSON(http.StatusCreated, createdDocument)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
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

	if c.GetHeader("x-ms-cosmos-is-query-plan-request") != "" {
		c.IndentedJSON(http.StatusOK, constants.QueryPlanResponse)
		return
	}

	var queryParameters map[string]interface{}
	if paramsArray, ok := requestBody["parameters"].([]interface{}); ok {
		queryParameters = parametersToMap(paramsArray)
	}

	docs, status := h.repository.ExecuteQueryDocuments(databaseId, collectionId, requestBody["query"].(string), queryParameters)
	if status != repositorymodels.StatusOk {
		// TODO: Currently we return everything if the query fails
		h.GetAllDocuments(c)
		return
	}

	collection, _ := h.repository.GetCollection(databaseId, collectionId)
	c.Header("x-ms-item-count", fmt.Sprintf("%d", len(docs)))
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
			createdDocument, status := h.repository.CreateDocument(databaseId, collectionId, operation.ResourceBody)
			responseCode := repositoryStatusToResponseCode(status)
			if status == repositorymodels.StatusOk {
				responseCode = http.StatusCreated
			}
			batchOperationResults[idx] = apimodels.BatchOperationResult{
				StatusCode:   responseCode,
				ResourceBody: createdDocument,
			}
		case apimodels.BatchOperationTypeDelete:
			status := h.repository.DeleteDocument(databaseId, collectionId, operation.Id)
			responseCode := repositoryStatusToResponseCode(status)
			if status == repositorymodels.StatusOk {
				responseCode = http.StatusNoContent
			}
			batchOperationResults[idx] = apimodels.BatchOperationResult{
				StatusCode: responseCode,
			}
		case apimodels.BatchOperationTypeReplace:
			deleteStatus := h.repository.DeleteDocument(databaseId, collectionId, operation.Id)
			if deleteStatus == repositorymodels.StatusNotFound {
				batchOperationResults[idx] = apimodels.BatchOperationResult{
					StatusCode: http.StatusNotFound,
				}
				continue
			}
			createdDocument, createStatus := h.repository.CreateDocument(databaseId, collectionId, operation.ResourceBody)
			responseCode := repositoryStatusToResponseCode(createStatus)
			if createStatus == repositorymodels.StatusOk {
				responseCode = http.StatusCreated
			}
			batchOperationResults[idx] = apimodels.BatchOperationResult{
				StatusCode:   responseCode,
				ResourceBody: createdDocument,
			}
		case apimodels.BatchOperationTypeUpsert:
			documentId := operation.ResourceBody["id"].(string)
			h.repository.DeleteDocument(databaseId, collectionId, documentId)
			createdDocument, createStatus := h.repository.CreateDocument(databaseId, collectionId, operation.ResourceBody)
			responseCode := repositoryStatusToResponseCode(createStatus)
			if createStatus == repositorymodels.StatusOk {
				responseCode = http.StatusCreated
			}
			batchOperationResults[idx] = apimodels.BatchOperationResult{
				StatusCode:   responseCode,
				ResourceBody: createdDocument,
			}
		case apimodels.BatchOperationTypeRead:
			document, status := h.repository.GetDocument(databaseId, collectionId, operation.Id)
			batchOperationResults[idx] = apimodels.BatchOperationResult{
				StatusCode:   repositoryStatusToResponseCode(status),
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

func repositoryStatusToResponseCode(status repositorymodels.RepositoryStatus) int {
	switch status {
	case repositorymodels.StatusOk:
		return http.StatusOK
	case repositorymodels.StatusNotFound:
		return http.StatusNotFound
	case repositorymodels.Conflict:
		return http.StatusConflict
	case repositorymodels.BadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
