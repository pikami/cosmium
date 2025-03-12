package badgerdatastore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/logger"
	"github.com/pikami/cosmium/internal/resourceid"
)

func (r *BadgerDataStore) GetAllDocuments(databaseId string, collectionId string) ([]datastore.Document, datastore.DataStoreStatus) {
	txn := r.db.NewTransaction(false)
	defer txn.Discard()

	dbExists, err := keyExists(txn, generateDatabaseKey(databaseId))
	if err != nil || !dbExists {
		return nil, datastore.StatusNotFound
	}

	collExists, err := keyExists(txn, generateCollectionKey(databaseId, collectionId))
	if err != nil || !collExists {
		return nil, datastore.StatusNotFound
	}

	docs, status := listByPrefix[datastore.Document](r.db, generateKey(resourceid.ResourceTypeDocument, databaseId, collectionId, ""))
	if status == datastore.StatusOk {
		return docs, datastore.StatusOk
	}

	return nil, status
}

func (r *BadgerDataStore) GetDocumentIterator(databaseId string, collectionId string) (datastore.DocumentIterator, datastore.DataStoreStatus) {
	txn := r.db.NewTransaction(false)

	dbExists, err := keyExists(txn, generateDatabaseKey(databaseId))
	if err != nil || !dbExists {
		return nil, datastore.StatusNotFound
	}

	collExists, err := keyExists(txn, generateCollectionKey(databaseId, collectionId))
	if err != nil || !collExists {
		return nil, datastore.StatusNotFound
	}

	iter := NewBadgerDocumentIterator(txn, generateKey(resourceid.ResourceTypeDocument, databaseId, collectionId, ""))
	return iter, datastore.StatusOk
}

func (r *BadgerDataStore) GetDocument(databaseId string, collectionId string, documentId string) (datastore.Document, datastore.DataStoreStatus) {
	documentKey := generateDocumentKey(databaseId, collectionId, documentId)

	txn := r.db.NewTransaction(false)
	defer txn.Discard()

	var document datastore.Document
	status := getKey(txn, documentKey, &document)

	return document, status
}

func (r *BadgerDataStore) DeleteDocument(databaseId string, collectionId string, documentId string) datastore.DataStoreStatus {
	documentKey := generateDocumentKey(databaseId, collectionId, documentId)

	txn := r.db.NewTransaction(true)
	defer txn.Discard()

	exists, err := keyExists(txn, documentKey)
	if err != nil {
		return datastore.Unknown
	}
	if !exists {
		return datastore.StatusNotFound
	}

	err = txn.Delete([]byte(documentKey))
	if err != nil {
		logger.ErrorLn("Error while deleting document:", err)
		return datastore.Unknown
	}

	err = txn.Commit()
	if err != nil {
		logger.ErrorLn("Error while committing transaction:", err)
		return datastore.Unknown
	}

	return datastore.StatusOk
}

func (r *BadgerDataStore) CreateDocument(databaseId string, collectionId string, document map[string]interface{}) (datastore.Document, datastore.DataStoreStatus) {
	txn := r.db.NewTransaction(true)
	defer txn.Discard()

	var database datastore.Database
	status := getKey(txn, generateDatabaseKey(databaseId), &database)
	if status != datastore.StatusOk {
		return datastore.Document{}, status
	}

	var collection datastore.Collection
	status = getKey(txn, generateCollectionKey(databaseId, collectionId), &collection)
	if status != datastore.StatusOk {
		return datastore.Document{}, status
	}

	var ok bool
	var documentId string
	if documentId, ok = document["id"].(string); !ok || documentId == "" {
		documentId = fmt.Sprint(uuid.New())
		document["id"] = documentId
	}

	document["_ts"] = time.Now().Unix()
	document["_rid"] = resourceid.NewCombined(collection.ResourceID, resourceid.New(resourceid.ResourceTypeDocument))
	document["_etag"] = fmt.Sprintf("\"%s\"", uuid.New())
	document["_self"] = fmt.Sprintf("dbs/%s/colls/%s/docs/%s/", database.ResourceID, collection.ResourceID, document["_rid"])

	status = insertKey(txn, generateDocumentKey(databaseId, collectionId, documentId), document)
	if status != datastore.StatusOk {
		return datastore.Document{}, status
	}

	return document, datastore.StatusOk
}
