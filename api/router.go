package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/api/handlers"
	"github.com/pikami/cosmium/api/handlers/middleware"
)

func CreateRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.RequestLogger())

	router.GET("/dbs/:databaseId/colls/:collId/pkranges", handlers.GetPartitionKeyRanges)

	router.POST("/dbs/:databaseId/colls/:collId/docs", handlers.DocumentsPost)
	router.GET("/dbs/:databaseId/colls/:collId/docs", handlers.GetAllDocuments)
	router.GET("/dbs/:databaseId/colls/:collId/docs/:docId", handlers.GetDocument)
	router.DELETE("/dbs/:databaseId/colls/:collId/docs/:docId", handlers.DeleteDocument)

	router.POST("/dbs/:databaseId/colls", handlers.CreateCollection)
	router.GET("/dbs/:databaseId/colls", handlers.GetAllCollections)
	router.GET("/dbs/:databaseId/colls/:collId", handlers.GetCollection)
	router.DELETE("/dbs/:databaseId/colls/:collId", handlers.DeleteCollection)

	router.POST("/dbs", handlers.CreateDatabase)
	router.GET("/dbs", handlers.GetAllDatabases)
	router.GET("/dbs/:databaseId", handlers.GetDatabase)
	router.DELETE("/dbs/:databaseId", handlers.DeleteDatabase)

	router.GET("/dbs/:databaseId/colls/:collId/udfs", handlers.GetAllUserDefinedFunctions)
	router.GET("/dbs/:databaseId/colls/:collId/sprocs", handlers.GetAllStoredProcedures)
	router.GET("/dbs/:databaseId/colls/:collId/triggers", handlers.GetAllTriggers)

	router.GET("/offers", handlers.GetOffers)
	router.GET("/", handlers.GetServerInfo)

	return router
}
