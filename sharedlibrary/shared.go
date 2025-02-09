package main

import (
	"sync"

	"github.com/pikami/cosmium/api"
	"github.com/pikami/cosmium/internal/repositories"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

type ServerInstance struct {
	server     *api.ApiServer
	repository *repositories.DataRepository
}

var (
	serverInstances = make(map[string]*ServerInstance)
	mutex           = sync.Mutex{}
)

const (
	ResponseSuccess = 0

	ResponseUnknown                     = 100
	ResponseFailedToParseConfiguration  = 101
	ResponseFailedToLoadState           = 102
	ResponseFailedToParseRequest        = 103
	ResponseServerInstanceAlreadyExists = 104
	ResponseServerInstanceNotFound      = 105
	ResponseFailedToStartServer         = 106

	ResponseRepositoryNotFound   = 200
	ResponseRepositoryConflict   = 201
	ResponseRepositoryBadRequest = 202
)

func getInstance(serverName string) (*ServerInstance, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	var ok bool
	var serverInstance *ServerInstance
	if serverInstance, ok = serverInstances[serverName]; !ok {
		return nil, false
	}

	return serverInstance, true
}

func addInstance(serverName string, serverInstance *ServerInstance) {
	mutex.Lock()
	defer mutex.Unlock()

	serverInstances[serverName] = serverInstance
}

func removeInstance(serverName string) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(serverInstances, serverName)
}

func repositoryStatusToResponseCode(status repositorymodels.RepositoryStatus) int {
	switch status {
	case repositorymodels.StatusOk:
		return ResponseSuccess
	case repositorymodels.StatusNotFound:
		return ResponseRepositoryNotFound
	case repositorymodels.Conflict:
		return ResponseRepositoryConflict
	case repositorymodels.BadRequest:
		return ResponseRepositoryBadRequest
	default:
		return ResponseUnknown
	}
}
