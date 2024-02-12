package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/internal/repositories"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

func GetAllUserDefinedFunctions(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	udfs, status := repositories.GetAllUserDefinedFunctions(databaseId, collectionId)

	if status == repositorymodels.StatusOk {
		c.IndentedJSON(http.StatusOK, gin.H{"_rid": "", "UserDefinedFunctions": udfs, "_count": len(udfs)})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
