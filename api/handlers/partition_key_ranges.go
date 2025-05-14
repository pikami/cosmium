package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/internal/constants"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/resourceid"
)

func (h *Handlers) GetPartitionKeyRanges(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	if c.Request.Header.Get("if-none-match") != "" {
		c.AbortWithStatus(http.StatusNotModified)
		return
	}

	partitionKeyRanges, status := h.dataStore.GetPartitionKeyRanges(databaseId, collectionId)
	if status == datastore.StatusOk {
		c.Header("etag", "\"420\"")
		c.Header("lsn", "420")
		c.Header("x-ms-cosmos-llsn", "420")
		c.Header("x-ms-global-committed-lsn", "420")
		c.Header("x-ms-item-count", fmt.Sprintf("%d", len(partitionKeyRanges)))

		collectionRid := collectionId
		collection, _ := h.dataStore.GetCollection(databaseId, collectionId)
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

	if status == datastore.StatusNotFound {
		c.IndentedJSON(http.StatusNotFound, constants.NotFoundResponse)
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, constants.UnknownErrorResponse)
}
