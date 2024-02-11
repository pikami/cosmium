package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/internal/constants"
	"github.com/pikami/cosmium/internal/repositories"
)

func GetAllDocuments(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	documents, status := repositories.GetAllDocuments(databaseId, collectionId)
	if status == repositories.StatusOk {
		c.IndentedJSON(http.StatusOK, gin.H{"_rid": "", "Documents": documents, "_count": len(documents)})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func GetDocument(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	documentId := c.Param("docId")

	document, status := repositories.GetDocument(databaseId, collectionId, documentId)
	if status == repositories.StatusOk {
		c.IndentedJSON(http.StatusOK, document)
		return
	}

	if status == repositories.StatusNotFound {
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
	if status == repositories.StatusOk {
		c.Status(http.StatusNoContent)
		return
	}

	if status == repositories.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
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

		// TODO: Handle these {"query":"select c.id, c._self, c._rid, c._ts, [c[\"pk\"]] as _partitionKeyValue from c"}
		docs, status := repositories.ExecuteQueryDocuments(databaseId, collectionId, query.(string))
		if status != repositories.StatusOk {
			// TODO: Currently we return everything if the query fails
			GetAllDocuments(c)
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"_rid": "", "Documents": docs, "_count": len(docs)})
		return
	}

	if requestBody["id"] == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	status := repositories.CreateDocument(databaseId, collectionId, requestBody)
	if status == repositories.Conflict {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Conflict"})
		return
	}

	if status == repositories.StatusOk {
		c.IndentedJSON(http.StatusCreated, requestBody)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
