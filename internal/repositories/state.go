package repositories

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	repositorymodels "github.com/pikami/cosmium/internal/repository_models"
)

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
	fmt.Printf("Databases: %d\n", len(state.Databases))
	fmt.Printf("Collections: %d\n", len(state.Collections))
	fmt.Printf("Documents: %d\n", len(state.Documents))

	databases = state.Databases
	collections = state.Collections
	documents = state.Documents
}
