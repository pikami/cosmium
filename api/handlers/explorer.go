package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) RegisterExplorerHandlers(router *gin.Engine) {
	explorer := router.Group(h.config.ExplorerBaseUrlLocation)
	{
		explorer.Use(func(ctx *gin.Context) {
			if ctx.Param("filepath") == "/config.json" {
				endpoint := fmt.Sprintf("https://%s:%d", h.config.Host, h.config.Port)
				ctx.JSON(200, gin.H{
					"BACKEND_ENDPOINT":       endpoint,
					"MONGO_BACKEND_ENDPOINT": endpoint,
					"PROXY_PATH":             "/",
					"EMULATOR_ENDPOINT":      endpoint,
				})
				ctx.Abort()
			} else {
				ctx.Next()
			}
		})

		if h.config.ExplorerPath != "" {
			explorer.Static("/", h.config.ExplorerPath)
		}
	}
}
