package mapdatastore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/resourceid"
	"golang.org/x/exp/maps"
)

func (r *MapDataStore) GetAllDocuments(databaseId string, collectionId string) ([]datastore.Document, datastore.DataStoreStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return make([]datastore.Document, 0), datastore.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return make([]datastore.Document, 0), datastore.StatusNotFound
	}

	return maps.Values(r.storeState.Documents[databaseId][collectionId]), datastore.StatusOk
}

func (r *MapDataStore) GetDocument(databaseId string, collectionId string, documentId string) (datastore.Document, datastore.DataStoreStatus) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return datastore.Document{}, datastore.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return datastore.Document{}, datastore.StatusNotFound
	}

	if _, ok := r.storeState.Documents[databaseId][collectionId][documentId]; !ok {
		return datastore.Document{}, datastore.StatusNotFound
	}

	return r.storeState.Documents[databaseId][collectionId][documentId], datastore.StatusOk
}

func (r *MapDataStore) DeleteDocument(databaseId string, collectionId string, documentId string) datastore.DataStoreStatus {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	if _, ok := r.storeState.Databases[databaseId]; !ok {
		return datastore.StatusNotFound
	}

	if _, ok := r.storeState.Collections[databaseId][collectionId]; !ok {
		return datastore.StatusNotFound
	}

	if _, ok := r.storeState.Documents[databaseId][collectionId][documentId]; !ok {
		return datastore.StatusNotFound
	}

	delete(r.storeState.Documents[databaseId][collectionId], documentId)

	return datastore.StatusOk
}

func (r *MapDataStore) CreateDocument(databaseId string, collectionId string, document map[string]interface{}) (datastore.Document, datastore.DataStoreStatus) {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	var ok bool
	var documentId string
	var database datastore.Database
	var collection datastore.Collection
	if documentId, ok = document["id"].(string); !ok || documentId == "" {
		documentId = fmt.Sprint(uuid.New())
		document["id"] = documentId
	}

	if database, ok = r.storeState.Databases[databaseId]; !ok {
		return datastore.Document{}, datastore.StatusNotFound
	}

	if collection, ok = r.storeState.Collections[databaseId][collectionId]; !ok {
		return datastore.Document{}, datastore.StatusNotFound
	}

	if _, ok := r.storeState.Documents[databaseId][collectionId][documentId]; ok {
		return datastore.Document{}, datastore.Conflict
	}

	document["_ts"] = time.Now().Unix()
	document["_rid"] = resourceid.NewCombined(collection.ResourceID, resourceid.New(resourceid.ResourceTypeDocument))
	document["_etag"] = fmt.Sprintf("\"%s\"", uuid.New())
	document["_self"] = fmt.Sprintf("dbs/%s/colls/%s/docs/%s/", database.ResourceID, collection.ResourceID, document["_rid"])

	r.storeState.Documents[databaseId][collectionId][documentId] = document

	return document, datastore.StatusOk
}

func (r *MapDataStore) GetDocumentIterator(databaseId string, collectionId string) (datastore.DocumentIterator, datastore.DataStoreStatus) {
	documents, status := r.GetAllDocuments(databaseId, collectionId)
	if status != datastore.StatusOk {
		return nil, status
	}

	return &ArrayDocumentIterator{
		documents: documents,
		index:     -1,
	}, datastore.StatusOk
}
