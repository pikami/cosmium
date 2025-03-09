package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) CosmiumExport(c *gin.Context) {
	dataStoreState, err := h.dataStore.DumpToJson()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", []byte(dataStoreState))
}
