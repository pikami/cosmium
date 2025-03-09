package main

import "C"
import (
	"encoding/json"
	"strings"

	"github.com/pikami/cosmium/internal/datastore"
)

//export CreateDocument
func CreateDocument(serverName *C.char, databaseId *C.char, collectionId *C.char, documentJson *C.char) int {
	serverNameStr := C.GoString(serverName)
	databaseIdStr := C.GoString(databaseId)
	collectionIdStr := C.GoString(collectionId)
	documentStr := C.GoString(documentJson)

	var ok bool
	var serverInstance *ServerInstance
	if serverInstance, ok = getInstance(serverNameStr); !ok {
		return ResponseServerInstanceNotFound
	}

	var document datastore.Document
	err := json.NewDecoder(strings.NewReader(documentStr)).Decode(&document)
	if err != nil {
		return ResponseFailedToParseRequest
	}

	_, code := serverInstance.dataStore.CreateDocument(databaseIdStr, collectionIdStr, document)

	return dataStoreStatusToResponseCode(code)
}

//export GetDocument
func GetDocument(serverName *C.char, databaseId *C.char, collectionId *C.char, documentId *C.char) *C.char {
	serverNameStr := C.GoString(serverName)
	databaseIdStr := C.GoString(databaseId)
	collectionIdStr := C.GoString(collectionId)
	documentIdStr := C.GoString(documentId)

	var ok bool
	var serverInstance *ServerInstance
	if serverInstance, ok = getInstance(serverNameStr); !ok {
		return C.CString("")
	}

	document, code := serverInstance.dataStore.GetDocument(databaseIdStr, collectionIdStr, documentIdStr)
	if code != datastore.StatusOk {
		return C.CString("")
	}

	documentJson, err := json.Marshal(document)
	if err != nil {
		return C.CString("")
	}
	return C.CString(string(documentJson))
}

//export GetAllDocuments
func GetAllDocuments(serverName *C.char, databaseId *C.char, collectionId *C.char) *C.char {
	serverNameStr := C.GoString(serverName)
	databaseIdStr := C.GoString(databaseId)
	collectionIdStr := C.GoString(collectionId)

	var ok bool
	var serverInstance *ServerInstance
	if serverInstance, ok = getInstance(serverNameStr); !ok {
		return C.CString("")
	}

	documents, code := serverInstance.dataStore.GetAllDocuments(databaseIdStr, collectionIdStr)
	if code != datastore.StatusOk {
		return C.CString("")
	}

	documentsJson, err := json.Marshal(documents)
	if err != nil {
		return C.CString("")
	}
	return C.CString(string(documentsJson))
}

//export UpdateDocument
func UpdateDocument(serverName *C.char, databaseId *C.char, collectionId *C.char, documentId *C.char, documentJson *C.char) int {
	serverNameStr := C.GoString(serverName)
	databaseIdStr := C.GoString(databaseId)
	collectionIdStr := C.GoString(collectionId)
	documentIdStr := C.GoString(documentId)
	documentStr := C.GoString(documentJson)

	var ok bool
	var serverInstance *ServerInstance
	if serverInstance, ok = getInstance(serverNameStr); !ok {
		return ResponseServerInstanceNotFound
	}

	var document datastore.Document
	err := json.Unmarshal([]byte(documentStr), &document)
	if err != nil {
		return ResponseFailedToParseRequest
	}

	code := serverInstance.dataStore.DeleteDocument(databaseIdStr, collectionIdStr, documentIdStr)
	if code != datastore.StatusOk {
		return dataStoreStatusToResponseCode(code)
	}

	_, code = serverInstance.dataStore.CreateDocument(databaseIdStr, collectionIdStr, document)
	return dataStoreStatusToResponseCode(code)
}

//export DeleteDocument
func DeleteDocument(serverName *C.char, databaseId *C.char, collectionId *C.char, documentId *C.char) int {
	serverNameStr := C.GoString(serverName)
	databaseIdStr := C.GoString(databaseId)
	collectionIdStr := C.GoString(collectionId)
	documentIdStr := C.GoString(documentId)

	var ok bool
	var serverInstance *ServerInstance
	if serverInstance, ok = getInstance(serverNameStr); !ok {
		return ResponseServerInstanceNotFound
	}

	code := serverInstance.dataStore.DeleteDocument(databaseIdStr, collectionIdStr, documentIdStr)

	return dataStoreStatusToResponseCode(code)
}
