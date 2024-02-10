package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/internal/repositories"
)

func GetAllDatabases(c *gin.Context) {
	databases, status := repositories.GetAllDatabases()
	if status == repositories.StatusOk {
		c.IndentedJSON(http.StatusOK, gin.H{"_rid": "", "Databases": databases})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func GetDatabase(c *gin.Context) {
	id := c.Param("id")

	database, status := repositories.GetDatabase(id)
	if status == repositories.StatusOk {
		c.IndentedJSON(http.StatusOK, database)
		return
	}

	if status == repositories.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func DeleteDatabase(c *gin.Context) {
	id := c.Param("id")

	status := repositories.DeleteDatabase(id)
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

func CreateDatabase(c *gin.Context) {
	var newDatabase repositories.Database

	if err := c.BindJSON(&newDatabase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if newDatabase.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	status := repositories.CreateDatabase(newDatabase)
	if status == repositories.Conflict {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Conflict"})
		return
	}

	if status == repositories.StatusOk {
		c.IndentedJSON(http.StatusCreated, newDatabase)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
