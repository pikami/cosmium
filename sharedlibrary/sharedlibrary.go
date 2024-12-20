package main

import "C"
import (
	"encoding/json"

	"github.com/pikami/cosmium/api"
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/repositories"
)

//export CreateServerInstance
func CreateServerInstance(serverName *C.char, configurationJSON *C.char) int {
	configStr := C.GoString(configurationJSON)
	serverNameStr := C.GoString(serverName)

	if _, ok := getInstance(serverNameStr); ok {
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

	addInstance(serverNameStr, &ServerInstance{
		server:     server,
		repository: repository,
	})

	return ResponseSuccess
}

//export StopServerInstance
func StopServerInstance(serverName *C.char) int {
	serverNameStr := C.GoString(serverName)

	if serverInstance, ok := getInstance(serverNameStr); ok {
		serverInstance.server.Stop()
		removeInstance(serverNameStr)
		return ResponseSuccess
	}

	return ResponseServerInstanceNotFound
}

//export GetServerInstanceState
func GetServerInstanceState(serverName *C.char) *C.char {
	serverNameStr := C.GoString(serverName)

	if serverInstance, ok := getInstance(serverNameStr); ok {
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

	if serverInstance, ok := getInstance(serverNameStr); ok {
		err := serverInstance.repository.LoadStateJSON(stateJSONStr)
		if err != nil {
			return ResponseFailedToLoadState
		}
		return ResponseSuccess
	}

	return ResponseServerInstanceNotFound
}

func main() {}
