package repositories

var databases = []Database{
	{ID: "db1"},
	{ID: "db2"},
}

func GetAllDatabases() ([]Database, RepositoryStatus) {
	return databases, StatusOk
}

func GetDatabase(id string) (Database, RepositoryStatus) {
	for _, db := range databases {
		if db.ID == id {
			return db, StatusOk
		}
	}

	return Database{}, StatusNotFound
}

func DeleteDatabase(id string) RepositoryStatus {
	for index, db := range databases {
		if db.ID == id {
			databases = append(databases[:index], databases[index+1:]...)
			return StatusOk
		}
	}

	return StatusNotFound
}

func CreateDatabase(newDatabase Database) RepositoryStatus {
	for _, db := range databases {
		if db.ID == newDatabase.ID {
			return Conflict
		}
	}

	databases = append(databases, newDatabase)
	return StatusOk
}
