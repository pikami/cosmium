package repositories

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
	"github.com/pikami/cosmium/internal/resourceid"
	"github.com/pikami/cosmium/parsers"
	"github.com/pikami/cosmium/parsers/nosql"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
	"golang.org/x/exp/maps"
)

func (r *DataRepository) GetAllDocuments(databaseId string, collectionId string) ([]repositorymodels.Document, repositorymodels.RepositoryStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return make([]repositorymodels.Document, 0), repositorymodels.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return make([]repositorymodels.Document, 0), repositorymodels.StatusNotFound
	}

	return maps.Values(r.storeState.Documents[databaseId][collectionId]), repositorymodels.StatusOk
}

func (r *DataRepository) GetDocument(databaseId string, collectionId string, documentId string) (repositorymodels.Document, repositorymodels.RepositoryStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return repositorymodels.Document{}, repositorymodels.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return repositorymodels.Document{}, repositorymodels.StatusNotFound
	}

	if _, ok := r.storeState.Documents[databaseId][collectionId][documentId]; !ok {
		return repositorymodels.Document{}, repositorymodels.StatusNotFound
	}

	return r.storeState.Documents[databaseId][collectionId][documentId], repositorymodels.StatusOk
}

func (r *DataRepository) DeleteDocument(databaseId string, collectionId string, documentId string) repositorymodels.RepositoryStatus {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return repositorymodels.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return repositorymodels.StatusNotFound
	}

	if _, ok := r.storeState.Documents[databaseId][collectionId][documentId]; !ok {
		return repositorymodels.StatusNotFound
	}

	delete(r.storeState.Documents[databaseId][collectionId], documentId)

	return repositorymodels.StatusOk
}

func (r *DataRepository) CreateDocument(databaseId string, collectionId string, document map[string]interface{}) (repositorymodels.Document, repositorymodels.RepositoryStatus) {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	var ok bool
	var documentId string
	var database repositorymodels.Database
	var collection repositorymodels.Collection
	if documentId, ok = document["id"].(string); !ok || documentId == "" {
		documentId = fmt.Sprint(uuid.New())
		document["id"] = documentId
	}

	if database, ok = r.storeState.Databases[databaseId]; !ok {
		return repositorymodels.Document{}, repositorymodels.StatusNotFound
	}

	if collection, ok = r.storeState.Collections[databaseId][collectionId]; !ok {
		return repositorymodels.Document{}, repositorymodels.StatusNotFound
	}

	if _, ok := r.storeState.Documents[databaseId][collectionId][documentId]; ok {
		return repositorymodels.Document{}, repositorymodels.Conflict
	}

	document["_ts"] = time.Now().Unix()
	document["_rid"] = resourceid.NewCombined(database.ResourceID, collection.ResourceID, resourceid.New())
	document["_etag"] = fmt.Sprintf("\"%s\"", uuid.New())
	document["_self"] = fmt.Sprintf("dbs/%s/colls/%s/docs/%s/", database.ResourceID, collection.ResourceID, document["_rid"])

	r.storeState.Documents[databaseId][collectionId][documentId] = document

	return document, repositorymodels.StatusOk
}

func (r *DataRepository) ExecuteQueryDocuments(databaseId string, collectionId string, query string, queryParameters map[string]interface{}) ([]memoryexecutor.RowType, repositorymodels.RepositoryStatus) {
	parsedQuery, err := nosql.Parse("", []byte(query))
	if err != nil {
		log.Printf("Failed to parse query: %s\nerr: %v", query, err)
		return nil, repositorymodels.BadRequest
	}

	collectionDocuments, status := r.GetAllDocuments(databaseId, collectionId)
	if status != repositorymodels.StatusOk {
		return nil, status
	}

	covDocs := make([]memoryexecutor.RowType, 0)
	for _, doc := range collectionDocuments {
		covDocs = append(covDocs, map[string]interface{}(doc))
	}

	if typedQuery, ok := parsedQuery.(parsers.SelectStmt); ok {
		typedQuery.Parameters = queryParameters
		return memoryexecutor.ExecuteQuery(typedQuery, covDocs), repositorymodels.StatusOk
	}

	return nil, repositorymodels.BadRequest
}
