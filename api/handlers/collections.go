package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

func (h *Handlers) GetAllCollections(c *gin.Context) {
	databaseId := c.Param("databaseId")

	collections, status := h.repository.GetAllCollections(databaseId)
	if status == repositorymodels.StatusOk {
		database, _ := h.repository.GetDatabase(databaseId)

		c.Header("x-ms-item-count", fmt.Sprintf("%d", len(collections)))
		c.IndentedJSON(http.StatusOK, gin.H{
			"_rid":                database.ResourceID,
			"DocumentCollections": collections,
			"_count":              len(collections),
		})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) GetCollection(c *gin.Context) {
	databaseId := c.Param("databaseId")
	id := c.Param("collId")

	collection, status := h.repository.GetCollection(databaseId, id)
	if status == repositorymodels.StatusOk {
		c.IndentedJSON(http.StatusOK, collection)
		return
	}

	if status == repositorymodels.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) DeleteCollection(c *gin.Context) {
	databaseId := c.Param("databaseId")
	id := c.Param("collId")

	status := h.repository.DeleteCollection(databaseId, id)
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

func (h *Handlers) CreateCollection(c *gin.Context) {
	databaseId := c.Param("databaseId")
	var newCollection repositorymodels.Collection

	if err := c.BindJSON(&newCollection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if newCollection.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	createdCollection, status := h.repository.CreateCollection(databaseId, newCollection)
	if status == repositorymodels.Conflict {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Conflict"})
		return
	}

	if status == repositorymodels.StatusOk {
		c.IndentedJSON(http.StatusCreated, createdCollection)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
