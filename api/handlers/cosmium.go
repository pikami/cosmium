package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) CosmiumExport(c *gin.Context) {
	repositoryState, err := h.repository.GetState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", []byte(repositoryState))
}
