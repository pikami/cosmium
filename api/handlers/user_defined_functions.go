package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

func (h *Handlers) GetAllUserDefinedFunctions(c *gin.Context) {
	databaseId := c.Param("databaseId")
	collectionId := c.Param("collId")

	udfs, status := h.repository.GetAllUserDefinedFunctions(databaseId, collectionId)

	if status == repositorymodels.StatusOk {
		c.Header("x-ms-item-count", fmt.Sprintf("%d", len(udfs)))
		c.IndentedJSON(http.StatusOK, gin.H{"_rid": "", "UserDefinedFunctions": udfs, "_count": len(udfs)})
		return
	}

	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
}
