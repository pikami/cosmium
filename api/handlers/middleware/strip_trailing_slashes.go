package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/api/config"
)

func StripTrailingSlashes(r *gin.Engine, config *config.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if len(path) > 1 && path[len(path)-1] == '/' && !strings.Contains(path, config.ExplorerBaseUrlLocation) {
			c.Request.URL.Path = path[:len(path)-1]
			r.HandleContext(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
