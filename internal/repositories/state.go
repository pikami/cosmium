package repositories

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

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

func LoadStateFS(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading state JSON file: %v", err)
	}

	var state repositorymodels.State
	if err := json.Unmarshal(data, &state); err != nil {
		log.Fatalf("Error unmarshalling state JSON: %v", err)
	}

	fmt.Println("Loaded state:")
	fmt.Printf("Databases: %d\n", getLength(state.Databases))
	fmt.Printf("Collections: %d\n", getLength(state.Collections))
	fmt.Printf("Documents: %d\n", getLength(state.Documents))

	storeState = state

	ensureStoreStateNoNullReferences()
}

func SaveStateFS(filePath string) {
	data, err := json.MarshalIndent(storeState, "", "\t")
	if err != nil {
		fmt.Printf("Failed to save state: %v\n", err)
		return
	}

	os.WriteFile(filePath, data, os.ModePerm)

	fmt.Println("Saved state:")
	fmt.Printf("Databases: %d\n", getLength(storeState.Databases))
	fmt.Printf("Collections: %d\n", getLength(storeState.Collections))
	fmt.Printf("Documents: %d\n", getLength(storeState.Documents))
}

func GetState() repositorymodels.State {
	return storeState
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
