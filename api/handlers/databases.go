package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

func (h *Handlers) GetAllDatabases(c *gin.Context) {
	databases, status := h.repository.GetAllDatabases()
	if status == repositorymodels.StatusOk {
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

	database, status := h.repository.GetDatabase(id)
	if status == repositorymodels.StatusOk {
		c.IndentedJSON(http.StatusOK, database)
		return
	}

	if status == repositorymodels.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func (h *Handlers) DeleteDatabase(c *gin.Context) {
	id := c.Param("databaseId")

	status := h.repository.DeleteDatabase(id)
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

func (h *Handlers) CreateDatabase(c *gin.Context) {
	var newDatabase repositorymodels.Database

	if err := c.BindJSON(&newDatabase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if newDatabase.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	createdDatabase, status := h.repository.CreateDatabase(newDatabase)
	if status == repositorymodels.Conflict {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Conflict"})
		return
	}

	if status == repositorymodels.StatusOk {
		c.IndentedJSON(http.StatusCreated, createdDatabase)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
