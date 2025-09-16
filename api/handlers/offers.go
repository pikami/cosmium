package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/api/headers"
)

func GetOffers(c *gin.Context) {
	c.Header(headers.ItemCount, "0")
	c.IndentedJSON(http.StatusOK, gin.H{
		"_rid":   "",
		"_count": 0,
		"Offers": []interface{}{},
	})
}
