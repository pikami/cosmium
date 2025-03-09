package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/internal/datastore"
)

func (h *Handlers) GetAllUserDefinedFunctions(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	udfs, status := h.dataStore.GetAllUserDefinedFunctions(databaseId, collectionId)

	if status == datastore.StatusOk {
		c.Header("x-ms-item-count", fmt.Sprintf("%d", len(udfs)))
		c.IndentedJSON(http.StatusOK, gin.H{"_rid": "", "UserDefinedFunctions": udfs, "_count": len(udfs)})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) GetUserDefinedFunction(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	udfId := c.Param("udfId")

	udf, status := h.dataStore.GetUserDefinedFunction(databaseId, collectionId, udfId)

	if status == datastore.StatusOk {
		c.IndentedJSON(http.StatusOK, udf)
		return
	}

	if status == datastore.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) DeleteUserDefinedFunction(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	udfId := c.Param("udfId")

	status := h.dataStore.DeleteUserDefinedFunction(databaseId, collectionId, udfId)
	if status == datastore.StatusOk {
		c.Status(http.StatusNoContent)
		return
	}

	if status == datastore.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) ReplaceUserDefinedFunction(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	udfId := c.Param("udfId")

	var udf datastore.UserDefinedFunction
	if err := c.BindJSON(&udf); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid body"})
		return
	}

	status := h.dataStore.DeleteUserDefinedFunction(databaseId, collectionId, udfId)
	if status == datastore.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	createdUdf, status := h.dataStore.CreateUserDefinedFunction(databaseId, collectionId, udf)
	if status == datastore.Conflict {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Conflict"})
		return
	}

	if status == datastore.StatusOk {
		c.IndentedJSON(http.StatusOK, createdUdf)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) CreateUserDefinedFunction(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	var udf datastore.UserDefinedFunction
	if err := c.BindJSON(&udf); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid body"})
		return
	}

	createdUdf, status := h.dataStore.CreateUserDefinedFunction(databaseId, collectionId, udf)
	if status == datastore.Conflict {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Conflict"})
		return
	}

	if status == datastore.StatusOk {
		c.IndentedJSON(http.StatusCreated, createdUdf)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
