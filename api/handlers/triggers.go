package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

func (h *Handlers) GetAllTriggers(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	triggers, status := h.repository.GetAllTriggers(databaseId, collectionId)

	if status == repositorymodels.StatusOk {
		c.Header("x-ms-item-count", fmt.Sprintf("%d", len(triggers)))
		c.IndentedJSON(http.StatusOK, gin.H{"_rid": "", "Triggers": triggers, "_count": len(triggers)})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
