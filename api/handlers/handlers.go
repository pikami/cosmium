package handlers

import (
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/repositories"
)

type Handlers struct {
	repository *repositories.DataRepository
	config     *config.ServerConfig
}

func NewHandlers(dataRepository *repositories.DataRepository, config *config.ServerConfig) *Handlers {
	return &Handlers{
		repository: dataRepository,
		config:     config,
	}
}
