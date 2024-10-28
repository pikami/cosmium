package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/api/handlers"
	"github.com/pikami/cosmium/api/handlers/middleware"
	"github.com/pikami/cosmium/internal/logger"
	tlsprovider "github.com/pikami/cosmium/internal/tls_provider"
)

func CreateRouter() *gin.Engine {
	router := gin.Default(func(e *gin.Engine) {
		e.RedirectTrailingSlash = false
	})

	if config.Config.Debug {
		router.Use(middleware.RequestLogger())
	}

	router.Use(middleware.StripTrailingSlashes(router))
	router.Use(middleware.Authentication())

	router.GET("/dbs/:databaseId/colls/:collId/pkranges", handlers.GetPartitionKeyRanges)

	router.POST("/dbs/:databaseId/colls/:collId/docs", handlers.DocumentsPost)
	router.GET("/dbs/:databaseId/colls/:collId/docs", handlers.GetAllDocuments)
	router.GET("/dbs/:databaseId/colls/:collId/docs/:docId", handlers.GetDocument)
	router.PUT("/dbs/:databaseId/colls/:collId/docs/:docId", handlers.ReplaceDocument)
	router.PATCH("/dbs/:databaseId/colls/:collId/docs/:docId", handlers.PatchDocument)
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

	router.GET("/cosmium/export", handlers.CosmiumExport)

	handlers.RegisterExplorerHandlers(router)

	return router
}

func StartAPI() {
	if !config.Config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := CreateRouter()
	listenAddress := fmt.Sprintf(":%d", config.Config.Port)

	if config.Config.TLS_CertificatePath != "" && config.Config.TLS_CertificateKey != "" {
		err := router.RunTLS(
			listenAddress,
			config.Config.TLS_CertificatePath,
			config.Config.TLS_CertificateKey)
		if err != nil {
			logger.Error("Failed to start HTTPS server:", err)
		}

		return
	}

	if config.Config.DisableTls {
		router.Run(listenAddress)
	}

	tlsConfig := tlsprovider.GetDefaultTlsConfig()
	server := &http.Server{
		Addr:      listenAddress,
		Handler:   router.Handler(),
		TLSConfig: tlsConfig,
	}

	logger.Infof("Listening and serving HTTPS on %s\n", server.Addr)
	err := server.ListenAndServeTLS("", "")
	if err != nil {
		logger.Error("Failed to start HTTPS server:", err)
	}

	router.Run()
}
