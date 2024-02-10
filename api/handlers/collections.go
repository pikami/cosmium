package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/internal/repositories"
)

func GetAllCollections(c *gin.Context) {
	databaseId := c.Param("databaseId")

	collections, status := repositories.GetAllCollections(databaseId)
	if status == repositories.StatusOk {
		c.IndentedJSON(http.StatusOK, gin.H{"_rid": "", "DocumentCollections": collections})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func GetCollection(c *gin.Context) {
	databaseId := c.Param("databaseId")
	id := c.Param("collId")

	collection, status := repositories.GetCollection(databaseId, id)
	if status == repositories.StatusOk {
		c.IndentedJSON(http.StatusOK, collection)
		return
	}

	if status == repositories.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}

func DeleteCollection(c *gin.Context) {
	databaseId := c.Param("databaseId")
	id := c.Param("collId")

	status := repositories.DeleteCollection(databaseId, id)
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

func CreateCollection(c *gin.Context) {
	databaseId := c.Param("databaseId")
	var newCollection repositories.Collection

	if err := c.BindJSON(&newCollection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if newCollection.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	status := repositories.CreateCollection(databaseId, newCollection)
	if status == repositories.Conflict {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Conflict"})
		return
	}

	if status == repositories.StatusOk {
		c.IndentedJSON(http.StatusCreated, newCollection)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
