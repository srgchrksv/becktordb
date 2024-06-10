package main

import (
	"fmt"
	"log"

	"github.com/srgchrksv/becktordb/becktordb"
	"github.com/srgchrksv/becktordb/services"
	"github.com/srgchrksv/becktordb/utils"
)

// EXAMPLE HOW TO USE
func main() {
	// load env vars : API_KEYS etc uses godotenv.Load()
	err := utils.LoadEnv("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	directoriesForRetrieval := []string{"../../documents"}
	databaseName := "becktorGemini.db"

	// init vectorDB - becktordb - repository
	db := becktordb.NewVectorDB(databaseName)
	// init business layer - services
	services := services.NewBectorDBservice(db)
	// Init gemini
	gemini := services.Gemini.Init()
	// embedd all documents in the directoriesForRetrieval directory
	err = gemini.GeminiEmbeddingsPath(directoriesForRetrieval, 50)
	if err != nil {
		log.Fatal(err)
	}
	// Query vectorDB - becktordb for similarity
	results, err := gemini.GeminiEmbeddingsQuery("What going on today?", 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Top 2 results:", results)
}
