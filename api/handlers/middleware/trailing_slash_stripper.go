package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func TrailingSlashStripper() gin.HandlerFunc {
	return func(c *gin.Context) {
		if (len(c.Request.URL.Path)) > 1 { //dont strip root dir slash, path="/"
			var stripped_path = strings.TrimSuffix(c.Request.URL.Path, "/")
			c.Request.URL.Path = stripped_path
		}
		c.Next()
	}
}
