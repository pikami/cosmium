package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/api/headers"
	"github.com/pikami/cosmium/internal/constants"
	"github.com/pikami/cosmium/internal/datastore"
)

func (h *Handlers) GetAllUserDefinedFunctions(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	udfs, status := h.dataStore.GetAllUserDefinedFunctions(databaseId, collectionId)

	if status == datastore.StatusOk {
		c.Header(headers.ItemCount, fmt.Sprintf("%d", len(udfs)))
		c.IndentedJSON(http.StatusOK, gin.H{"_rid": "", "UserDefinedFunctions": udfs, "_count": len(udfs)})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
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
		c.IndentedJSON(http.StatusNotFound, constants.NotFoundResponse)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
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
		c.IndentedJSON(http.StatusNotFound, constants.NotFoundResponse)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}

func (h *Handlers) ReplaceUserDefinedFunction(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	udfId := c.Param("udfId")

	var udf datastore.UserDefinedFunction
	if err := c.BindJSON(&udf); err != nil {
		c.IndentedJSON(http.StatusBadRequest, constants.BadRequestResponse)
		return
	}

	status := h.dataStore.DeleteUserDefinedFunction(databaseId, collectionId, udfId)
	if status == datastore.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, constants.NotFoundResponse)
		return
	}

	createdUdf, status := h.dataStore.CreateUserDefinedFunction(databaseId, collectionId, udf)
	if status == datastore.Conflict {
		c.IndentedJSON(http.StatusConflict, constants.ConflictResponse)
		return
	}

	if status == datastore.StatusOk {
		c.IndentedJSON(http.StatusOK, createdUdf)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}

func (h *Handlers) CreateUserDefinedFunction(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	var udf datastore.UserDefinedFunction
	if err := c.BindJSON(&udf); err != nil {
		c.IndentedJSON(http.StatusBadRequest, constants.BadRequestResponse)
		return
	}

	createdUdf, status := h.dataStore.CreateUserDefinedFunction(databaseId, collectionId, udf)
	if status == datastore.Conflict {
		c.IndentedJSON(http.StatusConflict, constants.ConflictResponse)
		return
	}

	if status == datastore.StatusOk {
		c.IndentedJSON(http.StatusCreated, createdUdf)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}
