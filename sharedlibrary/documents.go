package main

import "C"
import (
	"encoding/json"

	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
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

	var document repositorymodels.Document
	err := json.Unmarshal([]byte(documentStr), &document)
	if err != nil {
		return ResponseFailedToParseRequest
	}

	_, code := serverInstance.repository.CreateDocument(databaseIdStr, collectionIdStr, document)

	return repositoryStatusToResponseCode(code)
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

	document, code := serverInstance.repository.GetDocument(databaseIdStr, collectionIdStr, documentIdStr)
	if code != repositorymodels.StatusOk {
		return C.CString("")
	}

	documentJson, _ := json.Marshal(document)
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

	documents, code := serverInstance.repository.GetAllDocuments(databaseIdStr, collectionIdStr)
	if code != repositorymodels.StatusOk {
		return C.CString("")
	}

	documentsJson, _ := json.Marshal(documents)
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

	var document repositorymodels.Document
	err := json.Unmarshal([]byte(documentStr), &document)
	if err != nil {
		return ResponseFailedToParseRequest
	}

	code := serverInstance.repository.DeleteDocument(databaseIdStr, collectionIdStr, documentIdStr)
	if code != repositorymodels.StatusOk {
		return repositoryStatusToResponseCode(code)
	}

	_, code = serverInstance.repository.CreateDocument(databaseIdStr, collectionIdStr, document)
	return repositoryStatusToResponseCode(code)
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

	code := serverInstance.repository.DeleteDocument(databaseIdStr, collectionIdStr, documentIdStr)

	return repositoryStatusToResponseCode(code)
}
