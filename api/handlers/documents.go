package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/internal/constants"
	"github.com/pikami/cosmium/internal/repositories"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

func GetAllDocuments(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	documents, status := repositories.GetAllDocuments(databaseId, collectionId)
	if status == repositorymodels.StatusOk {
		collection, _ := repositories.GetCollection(databaseId, collectionId)

		c.IndentedJSON(http.StatusOK, gin.H{
			"_rid":      collection.ID,
			"Documents": documents,
			"_count":    len(documents),
		})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func GetDocument(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	documentId := c.Param("docId")

	document, status := repositories.GetDocument(databaseId, collectionId, documentId)
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

func DeleteDocument(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	documentId := c.Param("docId")

	status := repositories.DeleteDocument(databaseId, collectionId, documentId)
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
func ReplaceDocument(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	documentId := c.Param("docId")

	var requestBody map[string]interface{}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	status := repositories.DeleteDocument(databaseId, collectionId, documentId)
	if status == repositorymodels.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	createdDocument, status := repositories.CreateDocument(databaseId, collectionId, requestBody)
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

func DocumentsPost(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	var requestBody map[string]interface{}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	query := requestBody["query"]
	if query != nil {
		if c.GetHeader("x-ms-cosmos-is-query-plan-request") != "" {
			c.IndentedJSON(http.StatusOK, constants.QueryPlanResponse)
			return
		}

		var queryParameters map[string]interface{}
		if paramsArray, ok := requestBody["parameters"].([]interface{}); ok {
			queryParameters = parametersToMap(paramsArray)
		}

		docs, status := repositories.ExecuteQueryDocuments(databaseId, collectionId, query.(string), queryParameters)
		if status != repositorymodels.StatusOk {
			// TODO: Currently we return everything if the query fails
			GetAllDocuments(c)
			return
		}

		collection, _ := repositories.GetCollection(databaseId, collectionId)
		c.IndentedJSON(http.StatusOK, gin.H{
			"_rid":      collection.ResourceID,
			"Documents": docs,
			"_count":    len(docs),
		})
		return
	}

	if requestBody["id"] == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	createdDocument, status := repositories.CreateDocument(databaseId, collectionId, requestBody)
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
