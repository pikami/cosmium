package jsondatastore

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"sync"

	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/logger"
)

type State struct {
	sync.RWMutex

	// Map databaseId -> Database
	Databases map[string]datastore.Database `json:"databases"`

	// Map databaseId -> collectionId -> Collection
	Collections map[string]map[string]datastore.Collection `json:"collections"`

	// Map databaseId -> collectionId -> documentId -> Documents
	Documents map[string]map[string]map[string]datastore.Document `json:"documents"`

	// Map databaseId -> collectionId -> triggerId -> Trigger
	Triggers map[string]map[string]map[string]datastore.Trigger `json:"triggers"`

	// Map databaseId -> collectionId -> spId -> StoredProcedure
	StoredProcedures map[string]map[string]map[string]datastore.StoredProcedure `json:"sprocs"`

	// Map databaseId -> collectionId -> udfId -> UserDefinedFunction
	UserDefinedFunctions map[string]map[string]map[string]datastore.UserDefinedFunction `json:"udfs"`
}

func (r *JsonDataStore) InitializeDataStore() {
	if r.initialDataFilePath != "" {
		r.LoadStateFS(r.initialDataFilePath)
		return
	}

	if r.persistDataFilePath != "" {
		stat, err := os.Stat(r.persistDataFilePath)
		if err != nil {
			return
		}

		if stat.IsDir() {
			logger.ErrorLn("Argument '-Persist' must be a path to file, not a directory.")
			os.Exit(1)
		}

		r.LoadStateFS(r.persistDataFilePath)
		return
	}
}

func (r *JsonDataStore) LoadStateFS(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading state JSON file: %v", err)
		return
	}

	err = r.LoadStateJSON(string(data))
	if err != nil {
		log.Fatalf("Error unmarshalling state JSON: %v", err)
	}
}

func (r *JsonDataStore) LoadStateJSON(jsonData string) error {
	r.storeState.Lock()
	defer r.storeState.Unlock()

	var state State
	if err := json.Unmarshal([]byte(jsonData), &state); err != nil {
		return err
	}

	r.storeState.Collections = state.Collections
	r.storeState.Databases = state.Databases
	r.storeState.Documents = state.Documents

	r.ensureStoreStateNoNullReferences()

	logger.InfoLn("Loaded state:")
	logger.Infof("Databases: %d\n", getLength(r.storeState.Databases))
	logger.Infof("Collections: %d\n", getLength(r.storeState.Collections))
	logger.Infof("Documents: %d\n", getLength(r.storeState.Documents))
	logger.Infof("Triggers: %d\n", getLength(r.storeState.Triggers))
	logger.Infof("Stored procedures: %d\n", getLength(r.storeState.StoredProcedures))
	logger.Infof("User defined functions: %d\n", getLength(r.storeState.UserDefinedFunctions))

	return nil
}

func (r *JsonDataStore) SaveStateFS(filePath string) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	data, err := json.MarshalIndent(r.storeState, "", "\t")
	if err != nil {
		logger.Errorf("Failed to save state: %v\n", err)
		return
	}

	os.WriteFile(filePath, data, os.ModePerm)

	logger.InfoLn("Saved state:")
	logger.Infof("Databases: %d\n", getLength(r.storeState.Databases))
	logger.Infof("Collections: %d\n", getLength(r.storeState.Collections))
	logger.Infof("Documents: %d\n", getLength(r.storeState.Documents))
	logger.Infof("Triggers: %d\n", getLength(r.storeState.Triggers))
	logger.Infof("Stored procedures: %d\n", getLength(r.storeState.StoredProcedures))
	logger.Infof("User defined functions: %d\n", getLength(r.storeState.UserDefinedFunctions))
}

func (r *JsonDataStore) DumpToJson() (string, error) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	data, err := json.MarshalIndent(r.storeState, "", "\t")
	if err != nil {
		logger.Errorf("Failed to serialize state: %v\n", err)
		return "", err
	}

	return string(data), nil

}

func (r *JsonDataStore) Close() {
	if r.persistDataFilePath != "" {
		r.SaveStateFS(r.persistDataFilePath)
	}
}

func getLength(v interface{}) int {
	switch v.(type) {
	case datastore.Database,
		datastore.Collection,
		datastore.Document,
		datastore.Trigger,
		datastore.StoredProcedure,
		datastore.UserDefinedFunction:
		return 1
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Map {
		return -1
	}

	count := 0
	for _, key := range rv.MapKeys() {
		if rv.MapIndex(key).Kind() == reflect.Map {
			count += getLength(rv.MapIndex(key).Interface())
		} else {
			count++
		}
	}

	return count
}

func (r *JsonDataStore) ensureStoreStateNoNullReferences() {
	if r.storeState.Databases == nil {
		r.storeState.Databases = make(map[string]datastore.Database)
	}

	if r.storeState.Collections == nil {
		r.storeState.Collections = make(map[string]map[string]datastore.Collection)
	}

	if r.storeState.Documents == nil {
		r.storeState.Documents = make(map[string]map[string]map[string]datastore.Document)
	}

	if r.storeState.Triggers == nil {
		r.storeState.Triggers = make(map[string]map[string]map[string]datastore.Trigger)
	}

	if r.storeState.StoredProcedures == nil {
		r.storeState.StoredProcedures = make(map[string]map[string]map[string]datastore.StoredProcedure)
	}

	if r.storeState.UserDefinedFunctions == nil {
		r.storeState.UserDefinedFunctions = make(map[string]map[string]map[string]datastore.UserDefinedFunction)
	}

	for database := range r.storeState.Databases {
		if r.storeState.Collections[database] == nil {
			r.storeState.Collections[database] = make(map[string]datastore.Collection)
		}

		if r.storeState.Documents[database] == nil {
			r.storeState.Documents[database] = make(map[string]map[string]datastore.Document)
		}

		if r.storeState.Triggers[database] == nil {
			r.storeState.Triggers[database] = make(map[string]map[string]datastore.Trigger)
		}

		if r.storeState.StoredProcedures[database] == nil {
			r.storeState.StoredProcedures[database] = make(map[string]map[string]datastore.StoredProcedure)
		}

		if r.storeState.UserDefinedFunctions[database] == nil {
			r.storeState.UserDefinedFunctions[database] = make(map[string]map[string]datastore.UserDefinedFunction)
		}

		for collection := range r.storeState.Collections[database] {
			if r.storeState.Documents[database][collection] == nil {
				r.storeState.Documents[database][collection] = make(map[string]datastore.Document)
			}

			for document := range r.storeState.Documents[database][collection] {
				if r.storeState.Documents[database][collection][document] == nil {
					delete(r.storeState.Documents[database][collection], document)
				}
			}

			if r.storeState.Triggers[database][collection] == nil {
				r.storeState.Triggers[database][collection] = make(map[string]datastore.Trigger)
			}

			if r.storeState.StoredProcedures[database][collection] == nil {
				r.storeState.StoredProcedures[database][collection] = make(map[string]datastore.StoredProcedure)
			}

			if r.storeState.UserDefinedFunctions[database][collection] == nil {
				r.storeState.UserDefinedFunctions[database][collection] = make(map[string]datastore.UserDefinedFunction)
			}
		}
	}
}
