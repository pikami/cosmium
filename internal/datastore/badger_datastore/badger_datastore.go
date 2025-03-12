package badgerdatastore

import (
	"encoding/gob"

	"github.com/dgraph-io/badger/v4"
	"github.com/pikami/cosmium/internal/logger"
)

type BadgerDataStore struct {
	db *badger.DB
}

func NewBadgerDataStore() *BadgerDataStore {
	gob.Register([]interface{}{})

	badgerOpts := badger.DefaultOptions("").WithInMemory(true)

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
