package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/internal/repositories"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

func GetAllStoredProcedures(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	sps, status := repositories.GetAllStoredProcedures(databaseId, collectionId)

	if status == repositorymodels.StatusOk {
		c.IndentedJSON(http.StatusOK, gin.H{"_rid": "", "StoredProcedures": sps, "_count": len(sps)})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
