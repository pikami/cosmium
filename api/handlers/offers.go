package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOffers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"_rid":   "",
		"_count": 0,
		"Offers": []interface{}{},
	})
}
