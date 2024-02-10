package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/internal/repositories"
)

func GetAllUserDefinedFunctions(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	udfs, status := repositories.GetAllUserDefinedFunctions(databaseId, collectionId)

	if status == repositories.StatusOk {
		c.IndentedJSON(http.StatusOK, gin.H{"_rid": "", "UserDefinedFunctions": udfs, "_count": len(udfs)})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
