package mapdatastore

import "github.com/pikami/cosmium/internal/datastore"

type MapDataStore struct {
	storeState State

	initialDataFilePath string
	persistDataFilePath string
}

type MapDataStoreOptions struct {
	InitialDataFilePath string
	PersistDataFilePath string
}

func NewMapDataStore(options MapDataStoreOptions) *MapDataStore {
	dataStore := &MapDataStore{
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
