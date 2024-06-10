### becktordb - Vector Database

![The Go gopher was designed by Renee French. (http://reneefrench.blogspot.com/)](gopher.png)
*The Go gopher was  originally designed by Renee French. Source: https://golang.org/doc/gopher*
### Features:
- Persistant storage:
    - Writes and reads binary file
- Document loader with chunking
- Embeddings service:
    - Gemini 
        - To run get GEMINI_API_KEY at [https://ai.google.dev/](https://ai.google.dev/) and set at .env
    - Bedrock 
        - To run get access keys in IAM and set in .env file and enable 'cohere.embed-english' at Bedrock AWS:
            ```
            AWS_ACCESS_KEY_ID=
            AWS_REGION=
            AWS_SECRET_ACCESS_KEY=
            ```

- Vector similarity search:
    - Cosine

## Usage examples in  `examples` directory:

### Gemini:
`examples/gemini/main.go`

```GO
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

```

### Bedrock:
`examples/bedrock/main.go`

```Go
db := becktordb.NewVectorDB(databaseName)
	// init business layer - services
	services := services.NewBectorDBservice(db)
	// init bedrock
	bedrock := services.Bedrock.Init()
	err = bedrock.BedrockEmbeddingsPath(directoriesForRetrieval, 50)
	if err != nil {
		log.Fatal(err)
	}
    // Query vectorDB - becktordb for similarity
	results, err := bedrock.BedrockEmbeddingsQuery("What going on today?", 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Top 2 results:", results)

```

