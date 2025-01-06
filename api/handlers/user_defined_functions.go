package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

func (h *Handlers) GetAllUserDefinedFunctions(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	udfs, status := h.repository.GetAllUserDefinedFunctions(databaseId, collectionId)

	if status == repositorymodels.StatusOk {
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

	udf, status := h.repository.GetUserDefinedFunction(databaseId, collectionId, udfId)

	if status == repositorymodels.StatusOk {
		c.IndentedJSON(http.StatusOK, udf)
		return
	}

	if status == repositorymodels.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) DeleteUserDefinedFunction(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	udfId := c.Param("udfId")

	status := h.repository.DeleteUserDefinedFunction(databaseId, collectionId, udfId)
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

func (h *Handlers) ReplaceUserDefinedFunction(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	udfId := c.Param("udfId")

	var udf repositorymodels.UserDefinedFunction
	if err := c.BindJSON(&udf); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid body"})
		return
	}

	status := h.repository.DeleteUserDefinedFunction(databaseId, collectionId, udfId)
	if status == repositorymodels.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	createdUdf, status := h.repository.CreateUserDefinedFunction(databaseId, collectionId, udf)
	if status == repositorymodels.Conflict {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Conflict"})
		return
	}

	if status == repositorymodels.StatusOk {
		c.IndentedJSON(http.StatusOK, createdUdf)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) CreateUserDefinedFunction(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	var udf repositorymodels.UserDefinedFunction
	if err := c.BindJSON(&udf); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid body"})
		return
	}

	createdUdf, status := h.repository.CreateUserDefinedFunction(databaseId, collectionId, udf)
	if status == repositorymodels.Conflict {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Conflict"})
		return
	}

	if status == repositorymodels.StatusOk {
		c.IndentedJSON(http.StatusCreated, createdUdf)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
