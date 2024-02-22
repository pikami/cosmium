package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/internal/repositories"
)

func CosmiumExport(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, repositories.GetState())
}
