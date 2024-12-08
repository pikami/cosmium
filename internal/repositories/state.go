package repositories

import (
	"encoding/json"
	"log"
	"os"
	"reflect"

	"github.com/pikami/cosmium/api/config"
	"github.com/pikami/cosmium/internal/logger"
	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

var storedProcedures = []repositorymodels.StoredProcedure{}
var triggers = []repositorymodels.Trigger{}
var userDefinedFunctions = []repositorymodels.UserDefinedFunction{}
var storeState = repositorymodels.State{
	Databases:   make(map[string]repositorymodels.Database),
	Collections: make(map[string]map[string]repositorymodels.Collection),
	Documents:   make(map[string]map[string]map[string]repositorymodels.Document),
}

func InitializeRepository() {
	if config.Config.InitialDataFilePath != "" {
		LoadStateFS(config.Config.InitialDataFilePath)
		return
	}

	if config.Config.PersistDataFilePath != "" {
		stat, err := os.Stat(config.Config.PersistDataFilePath)
		if err != nil {
			return
		}

		if stat.IsDir() {
			logger.Error("Argument '-Persist' must be a path to file, not a directory.")
			os.Exit(1)
		}

		LoadStateFS(config.Config.PersistDataFilePath)
		return
	}
}

func LoadStateFS(filePath string) {
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

	storeState = state

	ensureStoreStateNoNullReferences()
}

func SaveStateFS(filePath string) {
	storeState.RLock()
	defer storeState.RUnlock()

	data, err := json.MarshalIndent(storeState, "", "\t")
	if err != nil {
		logger.Errorf("Failed to save state: %v\n", err)
		return
	}

	os.WriteFile(filePath, data, os.ModePerm)

	logger.Info("Saved state:")
	logger.Infof("Databases: %d\n", getLength(storeState.Databases))
	logger.Infof("Collections: %d\n", getLength(storeState.Collections))
	logger.Infof("Documents: %d\n", getLength(storeState.Documents))
}

func GetState() (string, error) {
	storeState.RLock()
	defer storeState.RUnlock()

	data, err := json.MarshalIndent(storeState, "", "\t")
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

func ensureStoreStateNoNullReferences() {
	if storeState.Databases == nil {
		storeState.Databases = make(map[string]repositorymodels.Database)
	}

	if storeState.Collections == nil {
		storeState.Collections = make(map[string]map[string]repositorymodels.Collection)
	}

	if storeState.Documents == nil {
		storeState.Documents = make(map[string]map[string]map[string]repositorymodels.Document)
	}

	for database := range storeState.Databases {
		if storeState.Collections[database] == nil {
			storeState.Collections[database] = make(map[string]repositorymodels.Collection)
		}

		if storeState.Documents[database] == nil {
			storeState.Documents[database] = make(map[string]map[string]repositorymodels.Document)
		}

		for collection := range storeState.Collections[database] {
			if storeState.Documents[database][collection] == nil {
				storeState.Documents[database][collection] = make(map[string]repositorymodels.Document)
			}

			for document := range storeState.Documents[database][collection] {
				if storeState.Documents[database][collection][document] == nil {
					delete(storeState.Documents[database][collection], document)
				}
			}
		}
	}
}
