package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/internal/repositories"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

func GetAllCollections(c *gin.Context) {
	databaseId := c.Param("databaseId")

	collections, status := repositories.GetAllCollections(databaseId)
	if status == repositorymodels.StatusOk {
		database, _ := repositories.GetDatabase(databaseId)

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

func GetCollection(c *gin.Context) {
	databaseId := c.Param("databaseId")
	id := c.Param("collId")

	collection, status := repositories.GetCollection(databaseId, id)
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

func DeleteCollection(c *gin.Context) {
	databaseId := c.Param("databaseId")
	id := c.Param("collId")

	status := repositories.DeleteCollection(databaseId, id)
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

func CreateCollection(c *gin.Context) {
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

	createdCollection, status := repositories.CreateCollection(databaseId, newCollection)
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
