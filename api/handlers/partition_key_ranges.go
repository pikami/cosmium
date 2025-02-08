package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
	"github.com/pikami/cosmium/internal/resourceid"
)

func (h *Handlers) GetPartitionKeyRanges(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	if c.Request.Header.Get("if-none-match") != "" {
		c.AbortWithStatus(http.StatusNotModified)
		return
	}

	partitionKeyRanges, status := h.repository.GetPartitionKeyRanges(databaseId, collectionId)
	if status == repositorymodels.StatusOk {
		c.Header("etag", "\"420\"")
		c.Header("lsn", "420")
		c.Header("x-ms-cosmos-llsn", "420")
		c.Header("x-ms-global-committed-lsn", "420")
		c.Header("x-ms-item-count", fmt.Sprintf("%d", len(partitionKeyRanges)))

		collectionRid := collectionId
		collection, _ := h.repository.GetCollection(databaseId, collectionId)
		if collection.ResourceID != "" {
			collectionRid = collection.ResourceID
		}

		rid := resourceid.NewCombined(collectionRid, resourceid.New(resourceid.ResourceTypePartitionKeyRange))
		c.IndentedJSON(http.StatusOK, gin.H{
			"_rid":               rid,
			"_count":             len(partitionKeyRanges),
			"PartitionKeyRanges": partitionKeyRanges,
		})
		return
	}

	if status == repositorymodels.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "NotFound"})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
