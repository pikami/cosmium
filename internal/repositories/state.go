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

	var state repositorymodels.State
	if err := json.Unmarshal(data, &state); err != nil {
		log.Fatalf("Error unmarshalling state JSON: %v", err)
		return
	}

	logger.Info("Loaded state:")
	logger.Infof("Databases: %d\n", getLength(state.Databases))
	logger.Infof("Collections: %d\n", getLength(state.Collections))
	logger.Infof("Documents: %d\n", getLength(state.Documents))

	r.storeState = state

	r.ensureStoreStateNoNullReferences()
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
		repositorymodels.Document:
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

	for database := range r.storeState.Databases {
		if r.storeState.Collections[database] == nil {
			r.storeState.Collections[database] = make(map[string]repositorymodels.Collection)
		}

		if r.storeState.Documents[database] == nil {
			r.storeState.Documents[database] = make(map[string]map[string]repositorymodels.Document)
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
		}
	}
}
