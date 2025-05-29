package badgerdatastore

import (
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/pikami/cosmium/internal/logger"
)

type BadgerDataStore struct {
	db       *badger.DB
	gcTicker *time.Ticker
}

type BadgerDataStoreOptions struct {
	PersistDataFilePath string
}

func NewBadgerDataStore(options BadgerDataStoreOptions) *BadgerDataStore {
	badgerOpts := badger.DefaultOptions(options.PersistDataFilePath)
	badgerOpts = badgerOpts.WithLogger(newBadgerLogger())
	if options.PersistDataFilePath == "" {
		badgerOpts = badgerOpts.WithInMemory(true)
	}

	db, err := badger.Open(badgerOpts)
	if err != nil {
		panic(err)
	}

	gcTicker := time.NewTicker(5 * time.Minute)

	ds := &BadgerDataStore{
		db:       db,
		gcTicker: gcTicker,
	}

	go ds.runGarbageCollector()

	return ds
}

func (r *BadgerDataStore) Close() {
	if r.gcTicker != nil {
		r.gcTicker.Stop()
		r.gcTicker = nil
	}

	r.db.Close()
	r.db = nil
}

func (r *BadgerDataStore) DumpToJson() (string, error) {
	logger.ErrorLn("Badger datastore does not support state export currently.")
	return "{}", nil
}

func (r *BadgerDataStore) runGarbageCollector() {
	for range r.gcTicker.C {
	again:
		err := r.db.RunValueLogGC(0.7)
		if err == nil {
			goto again
		}
	}
}
