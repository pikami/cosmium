package repositories

import (
	"encoding/json"
	"log"
	"os"
	"reflect"

	"github.com/pikami/cosmium/internal/logger"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

func (r *DataRepository) InitializeRepository() {
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
			logger.Error("Argument '-Persist' must be a path to file, not a directory.")
			os.Exit(1)
		}

		r.LoadStateFS(r.persistDataFilePath)
		return
	}
}

func (r *DataRepository) LoadStateFS(filePath string) {
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

func (r *DataRepository) LoadStateJSON(jsonData string) error {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	var state repositorymodels.State
	if err := json.Unmarshal([]byte(jsonData), &state); err != nil {
		return err
	}

	r.storeState.Collections = state.Collections
	r.storeState.Databases = state.Databases
	r.storeState.Documents = state.Documents

	r.ensureStoreStateNoNullReferences()

	logger.Info("Loaded state:")
	logger.Infof("Databases: %d\n", getLength(r.storeState.Databases))
	logger.Infof("Collections: %d\n", getLength(r.storeState.Collections))
	logger.Infof("Documents: %d\n", getLength(r.storeState.Documents))
	logger.Infof("Triggers: %d\n", getLength(r.storeState.Triggers))
	logger.Infof("Stored procedures: %d\n", getLength(r.storeState.StoredProcedures))
	logger.Infof("User defined functions: %d\n", getLength(r.storeState.UserDefinedFunctions))

	return nil
}

func (r *DataRepository) SaveStateFS(filePath string) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	data, err := json.MarshalIndent(r.storeState, "", "\t")
	if err != nil {
		logger.Errorf("Failed to save state: %v\n", err)
		return
	}

	os.WriteFile(filePath, data, os.ModePerm)

	logger.Info("Saved state:")
	logger.Infof("Databases: %d\n", getLength(r.storeState.Databases))
	logger.Infof("Collections: %d\n", getLength(r.storeState.Collections))
	logger.Infof("Documents: %d\n", getLength(r.storeState.Documents))
	logger.Infof("Triggers: %d\n", getLength(r.storeState.Triggers))
	logger.Infof("Stored procedures: %d\n", getLength(r.storeState.StoredProcedures))
	logger.Infof("User defined functions: %d\n", getLength(r.storeState.UserDefinedFunctions))
}

func (r *DataRepository) GetState() (string, error) {
	r.storeState.RLock()
	defer r.storeState.RUnlock()

	data, err := json.MarshalIndent(r.storeState, "", "\t")
	if err != nil {
		logger.Errorf("Failed to serialize state: %v\n", err)
		return "", err
	}

	return string(data), nil
}

func getLength(v interface{}) int {
	switch v.(type) {
	case repositorymodels.Database,
		repositorymodels.Collection,
		repositorymodels.Document,
		repositorymodels.Trigger,
		repositorymodels.StoredProcedure,
		repositorymodels.UserDefinedFunction:
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

func (r *DataRepository) ensureStoreStateNoNullReferences() {
	if r.storeState.Databases == nil {
		r.storeState.Databases = make(map[string]repositorymodels.Database)
	}

	if r.storeState.Collections == nil {
		r.storeState.Collections = make(map[string]map[string]repositorymodels.Collection)
	}

	if r.storeState.Documents == nil {
		r.storeState.Documents = make(map[string]map[string]map[string]repositorymodels.Document)
	}

	if r.storeState.Triggers == nil {
		r.storeState.Triggers = make(map[string]map[string]map[string]repositorymodels.Trigger)
	}

	if r.storeState.StoredProcedures == nil {
		r.storeState.StoredProcedures = make(map[string]map[string]map[string]repositorymodels.StoredProcedure)
	}

	if r.storeState.UserDefinedFunctions == nil {
		r.storeState.UserDefinedFunctions = make(map[string]map[string]map[string]repositorymodels.UserDefinedFunction)
	}

	for database := range r.storeState.Databases {
		if r.storeState.Collections[database] == nil {
			r.storeState.Collections[database] = make(map[string]repositorymodels.Collection)
		}

		if r.storeState.Documents[database] == nil {
			r.storeState.Documents[database] = make(map[string]map[string]repositorymodels.Document)
		}

		if r.storeState.Triggers[database] == nil {
			r.storeState.Triggers[database] = make(map[string]map[string]repositorymodels.Trigger)
		}

		if r.storeState.StoredProcedures[database] == nil {
			r.storeState.StoredProcedures[database] = make(map[string]map[string]repositorymodels.StoredProcedure)
		}

		if r.storeState.UserDefinedFunctions[database] == nil {
			r.storeState.UserDefinedFunctions[database] = make(map[string]map[string]repositorymodels.UserDefinedFunction)
		}

		for collection := range r.storeState.Collections[database] {
			if r.storeState.Documents[database][collection] == nil {
				r.storeState.Documents[database][collection] = make(map[string]repositorymodels.Document)
			}

			for document := range r.storeState.Documents[database][collection] {
				if r.storeState.Documents[database][collection][document] == nil {
					delete(r.storeState.Documents[database][collection], document)
				}
			}

			if r.storeState.Triggers[database][collection] == nil {
				r.storeState.Triggers[database][collection] = make(map[string]repositorymodels.Trigger)
			}

			if r.storeState.StoredProcedures[database][collection] == nil {
				r.storeState.StoredProcedures[database][collection] = make(map[string]repositorymodels.StoredProcedure)
			}

			if r.storeState.UserDefinedFunctions[database][collection] == nil {
				r.storeState.UserDefinedFunctions[database][collection] = make(map[string]repositorymodels.UserDefinedFunction)
			}
		}
	}
}
