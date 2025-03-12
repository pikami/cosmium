package main

import (
	"sync"

	"github.com/pikami/cosmium/api"
	"github.com/pikami/cosmium/internal/datastore"
)

type ServerInstance struct {
	server    *api.ApiServer
	dataStore datastore.DataStore
}

var (
	serverInstances = make(map[string]*ServerInstance)
	mutex           = sync.Mutex{}
)

const (
	ResponseSuccess = 0

	ResponseUnknown                                    = 100
	ResponseFailedToParseConfiguration                 = 101
	ResponseFailedToLoadState                          = 102
	ResponseFailedToParseRequest                       = 103
	ResponseServerInstanceAlreadyExists                = 104
	ResponseServerInstanceNotFound                     = 105
	ResponseFailedToStartServer                        = 106
	ResponseCurrentDataStoreDoesNotSupportStateLoading = 107

	ResponseDataStoreNotFound   = 200
	ResponseDataStoreConflict   = 201
	ResponseDataStoreBadRequest = 202
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

func dataStoreStatusToResponseCode(status datastore.DataStoreStatus) int {
	switch status {
	case datastore.StatusOk:
		return ResponseSuccess
	case datastore.StatusNotFound:
		return ResponseDataStoreNotFound
	case datastore.Conflict:
		return ResponseDataStoreConflict
	case datastore.BadRequest:
		return ResponseDataStoreBadRequest
	default:
		return ResponseUnknown
	}
}
