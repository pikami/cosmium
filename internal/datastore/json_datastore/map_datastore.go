package jsondatastore

import "github.com/pikami/cosmium/internal/datastore"

type JsonDataStore struct {
	storeState State

	initialDataFilePath string
	persistDataFilePath string
}

type JsonDataStoreOptions struct {
	InitialDataFilePath string
	PersistDataFilePath string
}

func NewJsonDataStore(options JsonDataStoreOptions) *JsonDataStore {
	dataStore := &JsonDataStore{
		storeState: State{
			Databases:            make(map[string]datastore.Database),
			Collections:          make(map[string]map[string]datastore.Collection),
			Documents:            make(map[string]map[string]map[string]datastore.Document),
			Triggers:             make(map[string]map[string]map[string]datastore.Trigger),
			StoredProcedures:     make(map[string]map[string]map[string]datastore.StoredProcedure),
			UserDefinedFunctions: make(map[string]map[string]map[string]datastore.UserDefinedFunction),
		},
		initialDataFilePath: options.InitialDataFilePath,
		persistDataFilePath: options.PersistDataFilePath,
	}

	dataStore.InitializeDataStore()

	return dataStore
}
