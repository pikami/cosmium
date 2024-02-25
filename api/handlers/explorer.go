package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/api/config"
)

func RegisterExplorerHandlers(router *gin.Engine) {
	explorer := router.Group("/_explorer")
	{
		explorer.Use(func(ctx *gin.Context) {
			if ctx.Param("filepath") == "/config.json" {
				endpoint := fmt.Sprintf("https://%s:%d", config.Config.Host, config.Config.Port)
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

		if config.Config.ExplorerPath != "" {
			explorer.Static("/", config.Config.ExplorerPath)
		}
	}
}
