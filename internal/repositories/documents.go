package repositories

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
	"github.com/pikami/cosmium/parsers"
	"github.com/pikami/cosmium/parsers/nosql"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
)

var documents = []repositorymodels.Document{}

func GetAllDocuments(databaseId string, collectionId string) ([]repositorymodels.Document, repositorymodels.RepositoryStatus) {
	filteredDocuments := make([]repositorymodels.Document, 0)

	for _, doc := range documents {
		docDbId := doc["_internal"].(map[string]interface{})["databaseId"]
		docCollId := doc["_internal"].(map[string]interface{})["collectionId"]

		if docDbId == databaseId && docCollId == collectionId {
			doc["_partitionKeyValue"] = doc["_internal"].(map[string]interface{})["partitionKeyValue"]
			filteredDocuments = append(filteredDocuments, doc)
		}
	}

	return filteredDocuments, repositorymodels.StatusOk
}

func GetDocument(databaseId string, collectionId string, documentId string) (repositorymodels.Document, repositorymodels.RepositoryStatus) {
	for _, doc := range documents {
		docDbId := doc["_internal"].(map[string]interface{})["databaseId"]
		docCollId := doc["_internal"].(map[string]interface{})["collectionId"]
		docId := doc["id"]

		if docDbId == databaseId && docCollId == collectionId && docId == documentId {
			doc["_partitionKeyValue"] = doc["_internal"].(map[string]interface{})["partitionKeyValue"]
			return doc, repositorymodels.StatusOk
		}
	}

	return repositorymodels.Document{}, repositorymodels.StatusNotFound
}

func DeleteDocument(databaseId string, collectionId string, documentId string) repositorymodels.RepositoryStatus {
	for index, doc := range documents {
		docDbId := doc["_internal"].(map[string]interface{})["databaseId"]
		docCollId := doc["_internal"].(map[string]interface{})["collectionId"]
		docId := doc["id"]

		if docDbId == databaseId && docCollId == collectionId && docId == documentId {
			documents = append(documents[:index], documents[index+1:]...)
			return repositorymodels.StatusOk
		}
	}

	return repositorymodels.StatusNotFound
}

func CreateDocument(databaseId string, collectionId string, document map[string]interface{}) repositorymodels.RepositoryStatus {
	if document["id"] == "" {
		return repositorymodels.BadRequest
	}

	collection, status := GetCollection(databaseId, collectionId)
	if status != repositorymodels.StatusOk {
		return repositorymodels.StatusNotFound
	}

	for _, doc := range documents {
		docDbId := doc["_internal"].(map[string]interface{})["databaseId"]
		docCollId := doc["_internal"].(map[string]interface{})["collectionId"]
		docId := doc["id"]

		if docDbId == databaseId && docCollId == collectionId && docId == document["id"] {
			return repositorymodels.Conflict
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

	document["_ts"] = time.Now().Unix()
	document["_rid"] = uuid.New().String()
	document["_etag"] = fmt.Sprintf("\"%s\"", document["_rid"])
	document["_internal"] = map[string]interface{}{
		"databaseId":        databaseId,
		"collectionId":      collectionId,
		"partitionKeyValue": partitionKeyValue,
	}
	documents = append(documents, document)

	return repositorymodels.StatusOk
}

func ExecuteQueryDocuments(databaseId string, collectionId string, query string, queryParameters map[string]interface{}) ([]memoryexecutor.RowType, repositorymodels.RepositoryStatus) {
	parsedQuery, err := nosql.Parse("", []byte(query))
	if err != nil {
		log.Printf("Failed to parse query: %s\nerr: %v", query, err)
		return nil, repositorymodels.BadRequest
	}

	collectionDocuments, status := GetAllDocuments(databaseId, collectionId)
	if status != repositorymodels.StatusOk {
		return nil, status
	}

	covDocs := make([]memoryexecutor.RowType, 0)
	for _, doc := range collectionDocuments {
		covDocs = append(covDocs, map[string]interface{}(doc))
	}

	if typedQuery, ok := parsedQuery.(parsers.SelectStmt); ok {
		typedQuery.Parameters = queryParameters
		return memoryexecutor.Execute(typedQuery, covDocs), repositorymodels.StatusOk
	}

	return nil, repositorymodels.BadRequest
}
