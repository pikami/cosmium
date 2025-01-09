package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/api/handlers"
	"github.com/pikami/cosmium/api/handlers/middleware"
	"github.com/pikami/cosmium/internal/logger"
	"github.com/pikami/cosmium/internal/repositories"
	tlsprovider "github.com/pikami/cosmium/internal/tls_provider"
)

func (s *ApiServer) CreateRouter(repository *repositories.DataRepository) {
	routeHandlers := handlers.NewHandlers(repository, s.config)

	gin.DefaultWriter = logger.InfoWriter()
	gin.DefaultErrorWriter = logger.ErrorWriter()

	if s.config.LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default(func(e *gin.Engine) {
		e.RedirectTrailingSlash = false
	})

	if s.config.LogLevel == "debug" {
		router.Use(middleware.RequestLogger())
	}

	router.Use(middleware.StripTrailingSlashes(router, s.config))
	router.Use(middleware.Authentication(s.config))

	router.GET("/dbs/:databaseId/colls/:collId/pkranges", routeHandlers.GetPartitionKeyRanges)

	router.POST("/dbs/:databaseId/colls/:collId/docs", routeHandlers.DocumentsPost)
	router.GET("/dbs/:databaseId/colls/:collId/docs", routeHandlers.GetAllDocuments)
	router.GET("/dbs/:databaseId/colls/:collId/docs/:docId", routeHandlers.GetDocument)
	router.PUT("/dbs/:databaseId/colls/:collId/docs/:docId", routeHandlers.ReplaceDocument)
	router.PATCH("/dbs/:databaseId/colls/:collId/docs/:docId", routeHandlers.PatchDocument)
	router.DELETE("/dbs/:databaseId/colls/:collId/docs/:docId", routeHandlers.DeleteDocument)

	router.POST("/dbs/:databaseId/colls", routeHandlers.CreateCollection)
	router.GET("/dbs/:databaseId/colls", routeHandlers.GetAllCollections)
	router.GET("/dbs/:databaseId/colls/:collId", routeHandlers.GetCollection)
	router.DELETE("/dbs/:databaseId/colls/:collId", routeHandlers.DeleteCollection)

	router.POST("/dbs", routeHandlers.CreateDatabase)
	router.GET("/dbs", routeHandlers.GetAllDatabases)
	router.GET("/dbs/:databaseId", routeHandlers.GetDatabase)
	router.DELETE("/dbs/:databaseId", routeHandlers.DeleteDatabase)

	router.POST("/dbs/:databaseId/colls/:collId/triggers", routeHandlers.CreateTrigger)
	router.GET("/dbs/:databaseId/colls/:collId/triggers", routeHandlers.GetAllTriggers)
	router.GET("/dbs/:databaseId/colls/:collId/triggers/:triggerId", routeHandlers.GetTrigger)
	router.PUT("/dbs/:databaseId/colls/:collId/triggers/:triggerId", routeHandlers.ReplaceTrigger)
	router.DELETE("/dbs/:databaseId/colls/:collId/triggers/:triggerId", routeHandlers.DeleteTrigger)

	router.POST("/dbs/:databaseId/colls/:collId/sprocs", routeHandlers.CreateStoredProcedure)
	router.GET("/dbs/:databaseId/colls/:collId/sprocs", routeHandlers.GetAllStoredProcedures)
	router.GET("/dbs/:databaseId/colls/:collId/sprocs/:sprocId", routeHandlers.GetStoredProcedure)
	router.PUT("/dbs/:databaseId/colls/:collId/sprocs/:sprocId", routeHandlers.ReplaceStoredProcedure)
	router.DELETE("/dbs/:databaseId/colls/:collId/sprocs/:sprocId", routeHandlers.DeleteStoredProcedure)

	router.POST("/dbs/:databaseId/colls/:collId/udfs", routeHandlers.CreateUserDefinedFunction)
	router.GET("/dbs/:databaseId/colls/:collId/udfs", routeHandlers.GetAllUserDefinedFunctions)
	router.GET("/dbs/:databaseId/colls/:collId/udfs/:udfId", routeHandlers.GetUserDefinedFunction)
	router.PUT("/dbs/:databaseId/colls/:collId/udfs/:udfId", routeHandlers.ReplaceUserDefinedFunction)
	router.DELETE("/dbs/:databaseId/colls/:collId/udfs/:udfId", routeHandlers.DeleteUserDefinedFunction)

	router.GET("/offers", handlers.GetOffers)
	router.GET("/", routeHandlers.GetServerInfo)

	router.GET("/cosmium/export", routeHandlers.CosmiumExport)

	routeHandlers.RegisterExplorerHandlers(router)

	s.router = router
}

func (s *ApiServer) Start() {
	listenAddress := fmt.Sprintf(":%d", s.config.Port)
	s.isActive = true

	server := &http.Server{
		Addr:    listenAddress,
		Handler: s.router.Handler(),
	}

	go func() {
		<-s.stopServer
		logger.InfoLn("Shutting down server...")
		err := server.Shutdown(context.TODO())
		if err != nil {
			logger.ErrorLn("Failed to shutdown server:", err)
		}
	}()

	go func() {
		if s.config.DisableTls {
			logger.Infof("Listening and serving HTTP on %s\n", server.Addr)
			err := server.ListenAndServe()
			if err != nil {
				logger.ErrorLn("Failed to start HTTP server:", err)
			}
			s.isActive = false
		} else if s.config.TLS_CertificatePath != "" && s.config.TLS_CertificateKey != "" {
			logger.Infof("Listening and serving HTTPS on %s\n", server.Addr)
			err := server.ListenAndServeTLS(
				s.config.TLS_CertificatePath,
				s.config.TLS_CertificateKey)
			if err != nil {
				logger.ErrorLn("Failed to start HTTPS server:", err)
			}
			s.isActive = false
		} else {
			tlsConfig := tlsprovider.GetDefaultTlsConfig()
			server.TLSConfig = tlsConfig

			logger.Infof("Listening and serving HTTPS on %s\n", server.Addr)
			err := server.ListenAndServeTLS("", "")
			if err != nil {
				logger.ErrorLn("Failed to start HTTPS server:", err)
			}
			s.isActive = false
		}
	}()
}
