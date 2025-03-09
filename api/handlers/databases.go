package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/internal/datastore"
)

func (h *Handlers) GetAllDatabases(c *gin.Context) {
	databases, status := h.dataStore.GetAllDatabases()
	if status == datastore.StatusOk {
		c.Header("x-ms-item-count", fmt.Sprintf("%d", len(databases)))
		c.IndentedJSON(http.StatusOK, gin.H{
			"_rid":      "",
			"Databases": databases,
			"_count":    len(databases),
		})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) GetDatabase(c *gin.Context) {
	id := c.Param("databaseId")

	database, status := h.dataStore.GetDatabase(id)
	if status == datastore.StatusOk {
		c.IndentedJSON(http.StatusOK, database)
		return
	}

	if status == datastore.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) DeleteDatabase(c *gin.Context) {
	id := c.Param("databaseId")

	status := h.dataStore.DeleteDatabase(id)
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

func (h *Handlers) CreateDatabase(c *gin.Context) {
	var newDatabase datastore.Database

	if err := c.BindJSON(&newDatabase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if newDatabase.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	createdDatabase, status := h.dataStore.CreateDatabase(newDatabase)
	if status == datastore.Conflict {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Conflict"})
		return
	}

	if status == datastore.StatusOk {
		c.IndentedJSON(http.StatusCreated, createdDatabase)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
