package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"strings"
	"unsafe"

	"github.com/pikami/cosmium/api"
	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/datastore"
	badgerdatastore "github.com/pikami/cosmium/internal/datastore/badger_datastore"
	jsondatastore "github.com/pikami/cosmium/internal/datastore/json_datastore"
)

//export CreateServerInstance
func CreateServerInstance(serverName *C.char, configurationJSON *C.char) int {
	serverNameStr := C.GoString(serverName)
	configStr := C.GoString(configurationJSON)

	if _, ok := getInstance(serverNameStr); ok {
		return ResponseServerInstanceAlreadyExists
	}

	var configuration config.ServerConfig
	err := json.NewDecoder(strings.NewReader(configStr)).Decode(&configuration)
	if err != nil {
		return ResponseFailedToParseConfiguration
	}

	configuration.ApplyDefaultsToEmptyFields()
	configuration.PopulateCalculatedFields()

	var dataStore datastore.DataStore
	switch configuration.DataStore {
	case config.DataStoreBadger:
		dataStore = badgerdatastore.NewBadgerDataStore(badgerdatastore.BadgerDataStoreOptions{
			InitialDataFilePath: configuration.InitialDataFilePath,
			PersistDataFilePath: configuration.PersistDataFilePath,
		})
	default:
		dataStore = jsondatastore.NewJsonDataStore(jsondatastore.JsonDataStoreOptions{
			InitialDataFilePath: configuration.InitialDataFilePath,
			PersistDataFilePath: configuration.PersistDataFilePath,
		})
	}

	server := api.NewApiServer(dataStore, &configuration)
	err = server.Start()
	if err != nil {
		return ResponseFailedToStartServer
	}

	addInstance(serverNameStr, &ServerInstance{
		server:    server,
		dataStore: dataStore,
	})

	return ResponseSuccess
}

//export StopServerInstance
func StopServerInstance(serverName *C.char) int {
	serverNameStr := C.GoString(serverName)

	if serverInstance, ok := getInstance(serverNameStr); ok {
		serverInstance.server.Stop()
		serverInstance.dataStore.Close()
		removeInstance(serverNameStr)
		return ResponseSuccess
	}

	return ResponseServerInstanceNotFound
}

//export GetServerInstanceState
func GetServerInstanceState(serverName *C.char) *C.char {
	serverNameStr := C.GoString(serverName)

	if serverInstance, ok := getInstance(serverNameStr); ok {
		stateJSON, err := serverInstance.dataStore.DumpToJson()
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
		if jsonDS, ok := serverInstance.dataStore.(*jsondatastore.JsonDataStore); ok {
			err := jsonDS.LoadStateJSON(stateJSONStr)
			if err != nil {
				return ResponseFailedToLoadState
			}
			return ResponseSuccess
		}
		return ResponseCurrentDataStoreDoesNotSupportStateLoading
	}

	return ResponseServerInstanceNotFound
}

//export FreeMemory
func FreeMemory(ptr *C.char) {
	C.free(unsafe.Pointer(ptr))
}

func main() {}
