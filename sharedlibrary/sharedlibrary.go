package main

import "C"
import (
	"encoding/json"

	"github.com/pikami/cosmium/api"
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/repositories"
)

type ServerInstance struct {
	server     *api.ApiServer
	repository *repositories.DataRepository
}

var serverInstances map[string]*ServerInstance

const (
	ResponseSuccess                     = 0
	ResponseUnknown                     = 1
	ResponseServerInstanceAlreadyExists = 2
	ResponseFailedToParseConfiguration  = 3
	ResponseServerInstanceNotFound      = 4
	ResponseFailedToLoadState           = 5
)

//export CreateServerInstance
func CreateServerInstance(serverName *C.char, configurationJSON *C.char) int {
	configStr := C.GoString(configurationJSON)
	serverNameStr := C.GoString(serverName)

	if serverInstances == nil {
		serverInstances = make(map[string]*ServerInstance)
	}

	if _, ok := serverInstances[serverNameStr]; ok {
		return ResponseServerInstanceAlreadyExists
	}

	var configuration config.ServerConfig
	err := json.Unmarshal([]byte(configStr), &configuration)
	if err != nil {
		return ResponseFailedToParseConfiguration
	}

	configuration.PopulateCalculatedFields()
	configuration.ApplyDefaultsToEmptyFields()

	repository := repositories.NewDataRepository(repositories.RepositoryOptions{
		InitialDataFilePath: configuration.InitialDataFilePath,
		PersistDataFilePath: configuration.PersistDataFilePath,
	})

	server := api.NewApiServer(repository, configuration)
	server.Start()

	serverInstances[serverNameStr] = &ServerInstance{
		server:     server,
		repository: repository,
	}

	return ResponseSuccess
}

//export StopServerInstance
func StopServerInstance(serverName *C.char) int {
	if serverInstance, ok := serverInstances[C.GoString(serverName)]; ok {
		serverInstance.server.Stop()
		delete(serverInstances, C.GoString(serverName))
		return ResponseSuccess
	}

	return ResponseServerInstanceNotFound
}

//export GetServerInstanceState
func GetServerInstanceState(serverName *C.char) *C.char {
	if serverInstance, ok := serverInstances[C.GoString(serverName)]; ok {
		stateJSON, err := serverInstance.repository.GetState()
		if err != nil {
			return nil
		}
		return C.CString(stateJSON)
	}

	return nil
}

//export LoadServerInstanceState
func LoadServerInstanceState(serverName *C.char, stateJSON *C.char) int {
	serverNameStr := C.GoString(serverName)
	stateJSONStr := C.GoString(stateJSON)

	if serverInstance, ok := serverInstances[serverNameStr]; ok {
		err := serverInstance.repository.LoadStateJSON(stateJSONStr)
		if err != nil {
			return ResponseFailedToLoadState
		}
		return ResponseSuccess
	}

	return ResponseServerInstanceNotFound
}

func main() {}
