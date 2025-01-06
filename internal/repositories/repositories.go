package repositories

import repositorymodels "github.com/pikami/cosmium/internal/repository_models"

type DataRepository struct {
	storeState repositorymodels.State

	initialDataFilePath string
	persistDataFilePath string
}

type RepositoryOptions struct {
	InitialDataFilePath string
	PersistDataFilePath string
}

func NewDataRepository(options RepositoryOptions) *DataRepository {
	repository := &DataRepository{
		storeState: repositorymodels.State{
			Databases:            make(map[string]repositorymodels.Database),
			Collections:          make(map[string]map[string]repositorymodels.Collection),
			Documents:            make(map[string]map[string]map[string]repositorymodels.Document),
			Triggers:             make(map[string]map[string]map[string]repositorymodels.Trigger),
			StoredProcedures:     make(map[string]map[string]map[string]repositorymodels.StoredProcedure),
			UserDefinedFunctions: make(map[string]map[string]map[string]repositorymodels.UserDefinedFunction),
		},
		initialDataFilePath: options.InitialDataFilePath,
		persistDataFilePath: options.PersistDataFilePath,
	}

	repository.InitializeRepository()

	return repository
}
