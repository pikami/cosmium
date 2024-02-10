package repositories

type Database struct {
	ID string `json:"id"`
}

type RepositoryStatus int

const (
	StatusOk       = 1
	StatusNotFound = 2
	Conflict       = 3
)
