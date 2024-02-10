package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/internal/constants"
)

func GetServerInfo(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, constants.ServerInfoResponse)
}
