package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/internal/repositories"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

func GetAllTriggers(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	triggers, status := repositories.GetAllTriggers(databaseId, collectionId)

	if status == repositorymodels.StatusOk {
		c.Header("x-ms-item-count", fmt.Sprintf("%d", len(triggers)))
		c.IndentedJSON(http.StatusOK, gin.H{"_rid": "", "Triggers": triggers, "_count": len(triggers)})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
