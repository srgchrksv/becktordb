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
	databaseName := "becktorBedrock.db"

	db := becktordb.NewVectorDB(databaseName)
	// init business layer - services
	services := services.NewBectorDBservice(db)
	// init bedrock
	bedrock := services.Bedrock.Init()
	err = bedrock.BedrockEmbeddingsPath(directoriesForRetrieval, 50)
	if err != nil {
		log.Fatal(err)
	}
	results, err := bedrock.BedrockEmbeddingsQuery("What going on today?", 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Top 2 results:", results)
}
