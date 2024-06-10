package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/srgchrksv/becktordb/becktordb"
	"github.com/srgchrksv/becktordb/utils"
	"google.golang.org/api/option"
)

type OpenAI struct {
	db *becktordb.VectorDB
}

// OpenAIEmbeddingsQuery embed query and run vector similarity, return closest 2 chunks
func (e *OpenAI) OpenAIEmbeddingsQuery(query string, topK int) ([]becktordb.VectorKey, error) {
	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("OPENAI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	// For embeddings, use the embedding-001 model
	em := client.EmbeddingModel("embedding-001")
	res, err := em.EmbedContent(ctx, genai.Text(query))

	if err != nil {
		panic(err)
	}
	response, err := e.db.Query(utils.Float32ToFloat64(res.Embedding.Values), topK)
	return response, err
}

func (e *OpenAI) OpenAIEmbeddingsPath(directories []string, chunkSize int) error {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("OPENAI_API_KEY")))
	if err != nil {
		return err
	}
	defer client.Close()
	em := client.EmbeddingModel("embedding-001")
	b := em.NewBatch()

	all_chunks, err := Chunking(directories, chunkSize)
	if err != nil {
		return err
	}
	for _, chunk := range all_chunks {
		b = b.AddContent(genai.Text(chunk.TextChunk))
	}
	res, err := em.BatchEmbedContents(ctx, b)
	if err != nil {
		return err
	}

	for i, chunk := range res.Embeddings {
		e.db.Add(all_chunks[i], utils.Float32ToFloat64(chunk.Values))
	}
	err = e.db.SaveToFile()

	return err
}

// Define the request payload structure
type EmbeddingRequest struct {
	Input          []string `json:"input"`
	Model          string   `json:"model"`
	EncodingFormat string   `json:"encoding_format"`
}

// Define the response structure (simplified for example)
type EmbeddingResponse struct {
	Data [][]float64 `json:"data"`
}

// Function to get embeddings from OpenAI
func GetEmbeddings(input []string, apiKey string) ([]float64, error) {
	url := "https://api.openai.com/v1/embeddings"

	// Create the request payload
	payload := EmbeddingRequest{
		Input:          input,
		Model:          "text-embedding-ada-002",
		EncodingFormat: "float",
	}

	// Marshal the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set the necessary headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-200 response: %s, %s", resp.Status, string(bodyBytes))
	}

	// Read and parse the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var embeddingResponse EmbeddingResponse
	if err := json.Unmarshal(bodyBytes, &embeddingResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// Assuming the response has at least one embedding
	if len(embeddingResponse.Data) == 0 {
		return nil, fmt.Errorf("no embedding data found in response")
	}

	return embeddingResponse.Data[0], nil
}
