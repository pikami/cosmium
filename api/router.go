package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/api/handlers"
)

func CreateRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/dbs/:databaseId/colls", handlers.CreateCollection)
	router.GET("/dbs/:databaseId/colls", handlers.GetAllCollections)
	router.GET("/dbs/:databaseId/colls/:collId", handlers.GetCollection)
	router.DELETE("/dbs/:databaseId/colls/:collId", handlers.DeleteCollection)

	router.POST("/dbs", handlers.CreateDatabase)
	router.GET("/dbs", handlers.GetAllDatabases)
	router.GET("/dbs/:databaseId", handlers.GetDatabase)
	router.DELETE("/dbs/:databaseId", handlers.DeleteDatabase)

	router.GET("/", handlers.GetServerInfo)

	return router
}
