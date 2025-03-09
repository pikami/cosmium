package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/datastore"
)

type ApiServer struct {
	stopServer       chan interface{}
	onServerShutdown chan interface{}
	isActive         bool
	router           *gin.Engine
	config           *config.ServerConfig
}

func NewApiServer(dataStore datastore.DataStore, config *config.ServerConfig) *ApiServer {
	stopChan := make(chan interface{})
	onServerShutdownChan := make(chan interface{})

	apiServer := &ApiServer{
		stopServer:       stopChan,
		onServerShutdown: onServerShutdownChan,
		config:           config,
	}

	apiServer.CreateRouter(dataStore)

	return apiServer
}

func (s *ApiServer) GetRouter() *gin.Engine {
	return s.router
}

func (s *ApiServer) Stop() {
	s.stopServer <- true
	<-s.onServerShutdown
}
