package main

import "C"
import (
	"encoding/json"

	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

//export CreateDatabase
func CreateDatabase(serverName *C.char, databaseJson *C.char) int {
	serverNameStr := C.GoString(serverName)
	databaseStr := C.GoString(databaseJson)

	var ok bool
	var serverInstance *ServerInstance
	if serverInstance, ok = getInstance(serverNameStr); !ok {
		return ResponseServerInstanceNotFound
	}

	var database repositorymodels.Database
	err := json.Unmarshal([]byte(databaseStr), &database)
	if err != nil {
		return ResponseFailedToParseRequest
	}

	_, code := serverInstance.repository.CreateDatabase(database)

	return repositoryStatusToResponseCode(code)
}

//export GetDatabase
func GetDatabase(serverName *C.char, databaseId *C.char) *C.char {
	serverNameStr := C.GoString(serverName)
	databaseIdStr := C.GoString(databaseId)

	var ok bool
	var serverInstance *ServerInstance
	if serverInstance, ok = getInstance(serverNameStr); !ok {
		return C.CString("")
	}

	database, code := serverInstance.repository.GetDatabase(databaseIdStr)
	if code != repositorymodels.StatusOk {
		return C.CString("")
	}

	databaseJson, err := json.Marshal(database)
	if err != nil {
		return C.CString("")
	}
	return C.CString(string(databaseJson))
}

//export GetAllDatabases
func GetAllDatabases(serverName *C.char) *C.char {
	serverNameStr := C.GoString(serverName)

	var ok bool
	var serverInstance *ServerInstance
	if serverInstance, ok = getInstance(serverNameStr); !ok {
		return C.CString("")
	}

	databases, code := serverInstance.repository.GetAllDatabases()
	if code != repositorymodels.StatusOk {
		return C.CString("")
	}

	databasesJson, err := json.Marshal(databases)
	if err != nil {
		return C.CString("")
	}

	return C.CString(string(databasesJson))
}

//export DeleteDatabase
func DeleteDatabase(serverName *C.char, databaseId *C.char) int {
	serverNameStr := C.GoString(serverName)
	databaseIdStr := C.GoString(databaseId)

	var ok bool
	var serverInstance *ServerInstance
	if serverInstance, ok = getInstance(serverNameStr); !ok {
		return ResponseServerInstanceNotFound
	}

	code := serverInstance.repository.DeleteDatabase(databaseIdStr)

	return repositoryStatusToResponseCode(code)
}
