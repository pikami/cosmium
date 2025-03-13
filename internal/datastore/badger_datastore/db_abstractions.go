package badgerdatastore

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/logger"
	"github.com/pikami/cosmium/internal/resourceid"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	DatabaseKeyPrefix            = "DB:"
	CollectionKeyPrefix          = "COL:"
	DocumentKeyPrefix            = "DOC:"
	TriggerKeyPrefix             = "TRG:"
	StoredProcedureKeyPrefix     = "SP:"
	UserDefinedFunctionKeyPrefix = "UDF:"
)

func generateKey(
	resourceType resourceid.ResourceType,
	databaseId string,
	collectionId string,
	resourceId string,
) string {
	result := ""

	switch resourceType {
	case resourceid.ResourceTypeDatabase:
		result += DatabaseKeyPrefix
	case resourceid.ResourceTypeCollection:
		result += CollectionKeyPrefix
	case resourceid.ResourceTypeDocument:
		result += DocumentKeyPrefix
	case resourceid.ResourceTypeTrigger:
		result += TriggerKeyPrefix
	case resourceid.ResourceTypeStoredProcedure:
		result += StoredProcedureKeyPrefix
	case resourceid.ResourceTypeUserDefinedFunction:
		result += UserDefinedFunctionKeyPrefix
	}

	if databaseId != "" {
		result += databaseId
	}

	if collectionId != "" {
		result += "/colls/" + collectionId
	}

	if resourceId != "" {
		result += "/" + resourceId
	}

	return result
}

func generateDatabaseKey(databaseId string) string {
	return generateKey(resourceid.ResourceTypeDatabase, databaseId, "", "")
}

func generateCollectionKey(databaseId string, collectionId string) string {
	return generateKey(resourceid.ResourceTypeCollection, databaseId, collectionId, "")
}

func generateDocumentKey(databaseId string, collectionId string, documentId string) string {
	return generateKey(resourceid.ResourceTypeDocument, databaseId, collectionId, documentId)
}

func generateTriggerKey(databaseId string, collectionId string, triggerId string) string {
	return generateKey(resourceid.ResourceTypeTrigger, databaseId, collectionId, triggerId)
}

func generateStoredProcedureKey(databaseId string, collectionId string, storedProcedureId string) string {
	return generateKey(resourceid.ResourceTypeStoredProcedure, databaseId, collectionId, storedProcedureId)
}

func generateUserDefinedFunctionKey(databaseId string, collectionId string, udfId string) string {
	return generateKey(resourceid.ResourceTypeUserDefinedFunction, databaseId, collectionId, udfId)
}

func insertKey(txn *badger.Txn, key string, value interface{}) datastore.DataStoreStatus {
	_, err := txn.Get([]byte(key))
	if err == nil {
		return datastore.Conflict
	}

	if err != badger.ErrKeyNotFound {
		logger.ErrorLn("Error while checking if key exists:", err)
		return datastore.Unknown
	}

	buf, err := msgpack.Marshal(value)
	if err != nil {
		logger.ErrorLn("Error while encoding value:", err)
		return datastore.Unknown
	}

	err = txn.Set([]byte(key), buf)
	if err != nil {
		logger.ErrorLn("Error while setting key:", err)
		return datastore.Unknown
	}

	err = txn.Commit()
	if err != nil {
		logger.ErrorLn("Error while committing transaction:", err)
		return datastore.Unknown
	}

	return datastore.StatusOk
}

func getKey(txn *badger.Txn, key string, value interface{}) datastore.DataStoreStatus {
	item, err := txn.Get([]byte(key))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return datastore.StatusNotFound
		}
		logger.ErrorLn("Error while getting key:", err)
		return datastore.Unknown
	}

	val, err := item.ValueCopy(nil)
	if err != nil {
		logger.ErrorLn("Error while copying value:", err)
		return datastore.Unknown
	}

	if value == nil {
		logger.ErrorLn("getKey called with nil value")
		return datastore.Unknown
	}

	err = msgpack.Unmarshal(val, &value)
	if err != nil {
		logger.ErrorLn("Error while decoding value:", err)
		return datastore.Unknown
	}

	return datastore.StatusOk
}

func keyExists(txn *badger.Txn, key string) (bool, error) {
	_, err := txn.Get([]byte(key))
	if err == nil {
		return true, nil
	}

	if err == badger.ErrKeyNotFound {
		return false, nil
	}

	return false, err
}

func listByPrefix[T any](db *badger.DB, prefix string) ([]T, datastore.DataStoreStatus) {
	results := make([]T, 0)

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte(prefix)
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			var entry T

			status := getKey(txn, string(item.Key()), &entry)
			if status != datastore.StatusOk {
				logger.ErrorLn("Failed to retrieve entry:", string(item.Key()))
				continue
			}

			results = append(results, entry)
		}
		return nil
	})

	if err != nil {
		logger.ErrorLn("Error while listing entries:", err)
		return nil, datastore.Unknown
	}

	return results, datastore.StatusOk
}

func deleteKeysByPrefix(txn *badger.Txn, prefix string) error {
	opts := badger.DefaultIteratorOptions
	opts.Prefix = []byte(prefix)
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		key := it.Item().KeyCopy(nil)
		if err := txn.Delete(key); err != nil {
			logger.ErrorLn("Failed to delete key:", string(key), "Error:", err)
			return err
		}
	}

	return nil
}
