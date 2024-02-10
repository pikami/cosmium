package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/api/handlers"
)

func CreateRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/dbs/:id", handlers.GetDatabase)
	router.DELETE("/dbs/:id", handlers.DeleteDatabase)
	router.GET("/dbs", handlers.GetAllDatabases)
	router.POST("/dbs", handlers.CreateDatabase)
	router.GET("/", handlers.GetServerInfo)

	return router
}
