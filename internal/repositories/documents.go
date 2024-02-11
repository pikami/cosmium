package repositories

import (
	"log"
	"strings"

	"github.com/pikami/cosmium/parsers"
	"github.com/pikami/cosmium/parsers/nosql"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
)

var documents = []Document{}

func GetAllDocuments(databaseId string, collectionId string) ([]Document, RepositoryStatus) {
	filteredDocuments := make([]Document, 0)

	for _, doc := range documents {
		docDbId := doc["_internal"].(map[string]interface{})["databaseId"]
		docCollId := doc["_internal"].(map[string]interface{})["collectionId"]

		if docDbId == databaseId && docCollId == collectionId {
			doc["_partitionKeyValue"] = doc["_internal"].(map[string]interface{})["partitionKeyValue"]
			filteredDocuments = append(filteredDocuments, doc)
		}
	}

	return filteredDocuments, StatusOk
}

func GetDocument(databaseId string, collectionId string, documentId string) (Document, RepositoryStatus) {
	for _, doc := range documents {
		docDbId := doc["_internal"].(map[string]interface{})["databaseId"]
		docCollId := doc["_internal"].(map[string]interface{})["collectionId"]
		docId := doc["id"]

		if docDbId == databaseId && docCollId == collectionId && docId == documentId {
			doc["_partitionKeyValue"] = doc["_internal"].(map[string]interface{})["partitionKeyValue"]
			return doc, StatusOk
		}
	}

	return Document{}, StatusNotFound
}

func DeleteDocument(databaseId string, collectionId string, documentId string) RepositoryStatus {
	for index, doc := range documents {
		docDbId := doc["_internal"].(map[string]interface{})["databaseId"]
		docCollId := doc["_internal"].(map[string]interface{})["collectionId"]
		docId := doc["id"]

		if docDbId == databaseId && docCollId == collectionId && docId == documentId {
			documents = append(documents[:index], documents[index+1:]...)
			return StatusOk
		}
	}

	return StatusNotFound
}

func CreateDocument(databaseId string, collectionId string, document map[string]interface{}) RepositoryStatus {
	if document["id"] == "" {
		return BadRequest
	}

	collection, status := GetCollection(databaseId, collectionId)
	if status != StatusOk {
		return StatusNotFound
	}

	for _, doc := range documents {
		docDbId := doc["_internal"].(map[string]interface{})["databaseId"]
		docCollId := doc["_internal"].(map[string]interface{})["collectionId"]
		docId := doc["id"]

		if docDbId == databaseId && docCollId == collectionId && docId == document["id"] {
			return Conflict
		}
	}

	partitionKeyValue := make([]string, 0)
	for _, path := range collection.PartitionKey.Paths {
		var val interface{}
		for _, part := range strings.Split(path, "/") {
			val = document[part]
		}

		if val == nil {
			val = ""
		}

		// TODO: handle non-string partition keys
		partitionKeyValue = append(partitionKeyValue, val.(string))
	}

	document["_internal"] = map[string]interface{}{
		"databaseId":        databaseId,
		"collectionId":      collectionId,
		"partitionKeyValue": partitionKeyValue,
	}
	documents = append(documents, document)

	return StatusOk
}

func ExecuteQueryDocuments(databaseId string, collectionId string, query string) ([]memoryexecutor.RowType, RepositoryStatus) {
	parsedQuery, err := nosql.Parse("", []byte(query))
	if err != nil {
		log.Printf("Failed to parse query: %s\nerr: %v", query, err)
		return nil, BadRequest
	}

	collectionDocuments, status := GetAllDocuments(databaseId, collectionId)
	if status != StatusOk {
		return nil, status
	}

	covDocs := make([]memoryexecutor.RowType, 0)
	for _, doc := range collectionDocuments {
		covDocs = append(covDocs, map[string]interface{}(doc))
	}

	return memoryexecutor.Execute(parsedQuery.(parsers.SelectStmt), covDocs), StatusOk
}
