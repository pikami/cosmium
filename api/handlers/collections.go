package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/api/headers"
	"github.com/pikami/cosmium/internal/constants"
	"github.com/pikami/cosmium/internal/datastore"
)

func (h *Handlers) GetAllCollections(c *gin.Context) {
	databaseId := c.Param("databaseId")

	collections, status := h.dataStore.GetAllCollections(databaseId)
	if status == datastore.StatusOk {
		database, _ := h.dataStore.GetDatabase(databaseId)

		c.Header(headers.ItemCount, fmt.Sprintf("%d", len(collections)))
		c.IndentedJSON(http.StatusOK, gin.H{
			"_rid":                database.ResourceID,
			"DocumentCollections": collections,
			"_count":              len(collections),
		})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}

func (h *Handlers) GetCollection(c *gin.Context) {
	databaseId := c.Param("databaseId")
	id := c.Param("collId")

	collection, status := h.dataStore.GetCollection(databaseId, id)
	if status == datastore.StatusOk {
		c.IndentedJSON(http.StatusOK, collection)
		return
	}

	if status == datastore.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, constants.NotFoundResponse)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}

func (h *Handlers) DeleteCollection(c *gin.Context) {
	databaseId := c.Param("databaseId")
	id := c.Param("collId")

	status := h.dataStore.DeleteCollection(databaseId, id)
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

func (h *Handlers) CreateCollection(c *gin.Context) {
	databaseId := c.Param("databaseId")
	var newCollection datastore.Collection

	if err := c.BindJSON(&newCollection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if newCollection.ID == "" {
		c.JSON(http.StatusBadRequest, constants.BadRequestResponse)
		return
	}

	createdCollection, status := h.dataStore.CreateCollection(databaseId, newCollection)
	if status == datastore.Conflict {
		c.IndentedJSON(http.StatusConflict, constants.ConflictResponse)
		return
	}

	if status == datastore.StatusOk {
		c.IndentedJSON(http.StatusCreated, createdCollection)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}
