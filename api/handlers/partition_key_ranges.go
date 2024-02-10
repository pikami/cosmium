package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/internal/repositories"
)

func GetPartitionKeyRanges(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	partitionKeyRanges, status := repositories.GetPartitionKeyRanges(databaseId, collectionId)
	if status == repositories.StatusOk {
		c.IndentedJSON(http.StatusOK, gin.H{
			"_rid":               "",
			"_count":             len(partitionKeyRanges),
			"PartitionKeyRanges": partitionKeyRanges,
		})
		return
	}

	if status == repositories.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
