package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

func (h *Handlers) GetAllStoredProcedures(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	sps, status := h.repository.GetAllStoredProcedures(databaseId, collectionId)

	if status == repositorymodels.StatusOk {
		c.Header("x-ms-item-count", fmt.Sprintf("%d", len(sps)))
		c.IndentedJSON(http.StatusOK, gin.H{"_rid": "", "StoredProcedures": sps, "_count": len(sps)})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
