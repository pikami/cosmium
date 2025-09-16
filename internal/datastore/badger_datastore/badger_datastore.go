package badgerdatastore

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/logger"
)

type BadgerDataStore struct {
	db       *badger.DB
	gcTicker *time.Ticker
}

type BadgerDataStoreOptions struct {
	InitialDataFilePath string
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

	ds.initializeDataStore(options.InitialDataFilePath)

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

func (r *BadgerDataStore) initializeDataStore(initialDataFilePath string) {
	if initialDataFilePath == "" {
		return
	}

	stat, err := os.Stat(initialDataFilePath)
	if err != nil {
		panic(err)
	}

	if stat.IsDir() {
		logger.ErrorLn("Argument '-Persist' must be a path to file, not a directory.")
		os.Exit(1)
	}

	jsonData, err := os.ReadFile(initialDataFilePath)
	if err != nil {
		log.Fatalf("Error reading state JSON file: %v", err)
		return
	}

	var state datastore.InitialDataModel
	if err := json.Unmarshal([]byte(jsonData), &state); err != nil {
		log.Fatalf("Error parsing state JSON file: %v", err)
		return
	}

	for dbName, dbModel := range state.Databases {
		r.CreateDatabase(dbModel)
		for colName, colModel := range state.Collections[dbName] {
			r.CreateCollection(dbName, colModel)
			for _, docModel := range state.Documents[dbName][colName] {
				r.CreateDocument(dbName, colName, docModel)
			}

			for _, triggerModel := range state.Triggers[dbName][colName] {
				r.CreateTrigger(dbName, colName, triggerModel)
			}

			for _, spModel := range state.StoredProcedures[dbName][colName] {
				r.CreateStoredProcedure(dbName, colName, spModel)
			}

			for _, udfModel := range state.UserDefinedFunctions[dbName][colName] {
				r.CreateUserDefinedFunction(dbName, colName, udfModel)
			}
		}
	}
}
