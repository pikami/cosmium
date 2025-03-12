package badgerdatastore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/logger"
	"github.com/pikami/cosmium/internal/resourceid"
)

func (r *BadgerDataStore) GetAllDatabases() ([]datastore.Database, datastore.DataStoreStatus) {
	dbs, status := listByPrefix[datastore.Database](r.db, DatabaseKeyPrefix)
	if status == datastore.StatusOk {
		return dbs, datastore.StatusOk
	}

	return nil, status
}

func (r *BadgerDataStore) GetDatabase(id string) (datastore.Database, datastore.DataStoreStatus) {
	databaseKey := generateDatabaseKey(id)

	txn := r.db.NewTransaction(false)
	defer txn.Discard()

	var database datastore.Database
	status := getKey(txn, databaseKey, &database)

	return database, status
}

func (r *BadgerDataStore) DeleteDatabase(id string) datastore.DataStoreStatus {
	databaseKey := generateDatabaseKey(id)

	txn := r.db.NewTransaction(true)
	defer txn.Discard()

	prefixes := []string{
		generateKey(resourceid.ResourceTypeCollection, id, "", ""),
		generateKey(resourceid.ResourceTypeDocument, id, "", ""),
		generateKey(resourceid.ResourceTypeTrigger, id, "", ""),
		generateKey(resourceid.ResourceTypeStoredProcedure, id, "", ""),
		generateKey(resourceid.ResourceTypeUserDefinedFunction, id, "", ""),
		databaseKey,
	}
	for _, prefix := range prefixes {
		if err := deleteKeysByPrefix(txn, prefix); err != nil {
			return datastore.Unknown
		}
	}

	err := txn.Commit()
	if err != nil {
		logger.ErrorLn("Error while committing transaction:", err)
		return datastore.Unknown
	}

	return datastore.StatusOk
}

func (r *BadgerDataStore) CreateDatabase(newDatabase datastore.Database) (datastore.Database, datastore.DataStoreStatus) {
	databaseKey := generateDatabaseKey(newDatabase.ID)

	txn := r.db.NewTransaction(true)
	defer txn.Discard()

	newDatabase.TimeStamp = time.Now().Unix()
	newDatabase.ResourceID = resourceid.New(resourceid.ResourceTypeDatabase)
	newDatabase.ETag = fmt.Sprintf("\"%s\"", uuid.New())
	newDatabase.Self = fmt.Sprintf("dbs/%s/", newDatabase.ResourceID)

	status := insertKey(txn, databaseKey, newDatabase)
	if status != datastore.StatusOk {
		return datastore.Database{}, status
	}

	return newDatabase, datastore.StatusOk
}
