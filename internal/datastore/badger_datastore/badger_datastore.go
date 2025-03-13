package badgerdatastore

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/pikami/cosmium/internal/logger"
)

type BadgerDataStore struct {
	db *badger.DB
}

type BadgerDataStoreOptions struct {
	PersistDataFilePath string
}

func NewBadgerDataStore(options BadgerDataStoreOptions) *BadgerDataStore {
	badgerOpts := badger.DefaultOptions(options.PersistDataFilePath)
	if options.PersistDataFilePath == "" {
		badgerOpts = badgerOpts.WithInMemory(true)
	}

	db, err := badger.Open(badgerOpts)
	if err != nil {
		panic(err)
	}

	return &BadgerDataStore{
		db: db,
	}
}

func (r *BadgerDataStore) Close() {
	r.db.Close()
	r.db = nil
}

func (r *BadgerDataStore) DumpToJson() (string, error) {
	logger.ErrorLn("Badger datastore does not support state export currently.")
	return "{}", nil
}
