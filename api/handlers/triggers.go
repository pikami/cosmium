package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/internal/repositories"
)

func GetAllTriggers(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	triggers, status := repositories.GetAllTriggers(databaseId, collectionId)

	if status == repositories.StatusOk {
		c.IndentedJSON(http.StatusOK, gin.H{"_rid": "", "Triggers": triggers, "_count": len(triggers)})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
