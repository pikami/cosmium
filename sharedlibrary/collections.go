package main

import "C"
import (
	"encoding/json"
	"strings"

	"github.com/pikami/cosmium/internal/datastore"
)

//export CreateCollection
func CreateCollection(serverName *C.char, databaseId *C.char, collectionJson *C.char) int {
	serverNameStr := C.GoString(serverName)
	databaseIdStr := C.GoString(databaseId)
	collectionStr := C.GoString(collectionJson)

	var ok bool
	var serverInstance *ServerInstance
	if serverInstance, ok = getInstance(serverNameStr); !ok {
		return ResponseServerInstanceNotFound
	}

	var collection datastore.Collection
	err := json.NewDecoder(strings.NewReader(collectionStr)).Decode(&collection)
	if err != nil {
		return ResponseFailedToParseRequest
	}

	_, code := serverInstance.dataStore.CreateCollection(databaseIdStr, collection)

	return dataStoreStatusToResponseCode(code)
}

//export GetCollection
func GetCollection(serverName *C.char, databaseId *C.char, collectionId *C.char) *C.char {
	serverNameStr := C.GoString(serverName)
	databaseIdStr := C.GoString(databaseId)
	collectionIdStr := C.GoString(collectionId)

	var ok bool
	var serverInstance *ServerInstance
	if serverInstance, ok = getInstance(serverNameStr); !ok {
		return C.CString("")
	}

	collection, code := serverInstance.dataStore.GetCollection(databaseIdStr, collectionIdStr)
	if code != datastore.StatusOk {
		return C.CString("")
	}

	collectionJson, err := json.Marshal(collection)
	if err != nil {
		return C.CString("")
	}
	return C.CString(string(collectionJson))
}

//export GetAllCollections
func GetAllCollections(serverName *C.char, databaseId *C.char) *C.char {
	serverNameStr := C.GoString(serverName)
	databaseIdStr := C.GoString(databaseId)

	var ok bool
	var serverInstance *ServerInstance
	if serverInstance, ok = getInstance(serverNameStr); !ok {
		return C.CString("")
	}

	collections, code := serverInstance.dataStore.GetAllCollections(databaseIdStr)
	if code != datastore.StatusOk {
		return C.CString("")
	}

	collectionsJson, err := json.Marshal(collections)
	if err != nil {
		return C.CString("")
	}
	return C.CString(string(collectionsJson))
}

//export DeleteCollection
func DeleteCollection(serverName *C.char, databaseId *C.char, collectionId *C.char) int {
	serverNameStr := C.GoString(serverName)
	databaseIdStr := C.GoString(databaseId)
	collectionIdStr := C.GoString(collectionId)

	var ok bool
	var serverInstance *ServerInstance
	if serverInstance, ok = getInstance(serverNameStr); !ok {
		return ResponseServerInstanceNotFound
	}

	code := serverInstance.dataStore.DeleteCollection(databaseIdStr, collectionIdStr)

	return dataStoreStatusToResponseCode(code)
}
