package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/internal/constants"
	"github.com/pikami/cosmium/internal/datastore"
)

func (h *Handlers) GetAllStoredProcedures(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	sps, status := h.dataStore.GetAllStoredProcedures(databaseId, collectionId)

	if status == datastore.StatusOk {
		c.Header("x-ms-item-count", fmt.Sprintf("%d", len(sps)))
		c.IndentedJSON(http.StatusOK, gin.H{"_rid": "", "StoredProcedures": sps, "_count": len(sps)})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}

func (h *Handlers) GetStoredProcedure(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	spId := c.Param("spId")

	sp, status := h.dataStore.GetStoredProcedure(databaseId, collectionId, spId)

	if status == datastore.StatusOk {
		c.IndentedJSON(http.StatusOK, sp)
		return
	}

	if status == datastore.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, constants.NotFoundResponse)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}

func (h *Handlers) DeleteStoredProcedure(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	spId := c.Param("spId")

	status := h.dataStore.DeleteStoredProcedure(databaseId, collectionId, spId)
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

func (h *Handlers) ReplaceStoredProcedure(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")
	spId := c.Param("spId")

	var sp datastore.StoredProcedure
	if err := c.BindJSON(&sp); err != nil {
		c.IndentedJSON(http.StatusBadRequest, constants.BadRequestResponse)
		return
	}

	status := h.dataStore.DeleteStoredProcedure(databaseId, collectionId, spId)
	if status == datastore.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, constants.NotFoundResponse)
		return
	}

	createdSP, status := h.dataStore.CreateStoredProcedure(databaseId, collectionId, sp)
	if status == datastore.Conflict {
		c.IndentedJSON(http.StatusConflict, constants.ConflictResponse)
		return
	}

	if status == datastore.StatusOk {
		c.IndentedJSON(http.StatusOK, createdSP)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}

func (h *Handlers) CreateStoredProcedure(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	var sp datastore.StoredProcedure
	if err := c.BindJSON(&sp); err != nil {
		c.IndentedJSON(http.StatusBadRequest, constants.BadRequestResponse)
		return
	}

	createdSP, status := h.dataStore.CreateStoredProcedure(databaseId, collectionId, sp)
	if status == datastore.Conflict {
		c.IndentedJSON(http.StatusConflict, constants.ConflictResponse)
		return
	}

	if status == datastore.StatusOk {
		c.IndentedJSON(http.StatusCreated, createdSP)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}
