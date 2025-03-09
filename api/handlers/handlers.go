package handlers

import (
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/datastore"
)

type Handlers struct {
	dataStore datastore.DataStore
	config    *config.ServerConfig
}

func NewHandlers(dataStore datastore.DataStore, config *config.ServerConfig) *Handlers {
	return &Handlers{
		dataStore: dataStore,
		config:    config,
	}
}
