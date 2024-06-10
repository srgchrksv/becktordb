package services

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/srgchrksv/becktordb/becktordb"
)

const (
	titanEmbeddingModelID = "cohere.embed-english-v3" //https://docs.aws.amazon.com/bedrock/latest/userguide/model-ids-arns.html
)

type Bedrock struct {
	db  *becktordb.VectorDB
	brc *bedrockruntime.Client
	mu  sync.RWMutex
}

func (b *Bedrock) Init() *Bedrock {
	// Create a new AWS session
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return &Bedrock{db: b.db, mu: sync.RWMutex{}}
	}

	// Create a new Bedrock service client
	brc := bedrockruntime.NewFromConfig(cfg)
	return &Bedrock{
		db:  b.db,
		brc: brc,
	}
}

type Request struct {
	Texts     []string `json:"texts"`
	InputType string   `json:"input_type"`
}

type Response struct {
	Id         string      `json:"id"`
	Embeddings [][]float64 `json:"embeddings"`
}

// BedrockEmbeddingsQuery function takes a query string and returns the embeddings for that query.
func (b *Bedrock) BedrockEmbeddingsQuery(query string, topK int) ([]becktordb.VectorKey, error) {
	resp, err := b.bedrockEmbeddingsRequest([]string{query})
	if err != nil {
		return []becktordb.VectorKey{}, err
	}
	response, err := b.db.Query(resp.Embeddings[0], topK)
	return response, err
}

func (b *Bedrock) BedrockEmbeddingsPath(directories []string, toEmbedChunksize int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	allChunks, err := Chunking(directories, toEmbedChunksize)
	if err != nil {
		return err
	}
	var toEmbedChunks []string
	for _, chunk := range allChunks {
		toEmbedChunks = append(toEmbedChunks, chunk.TextChunk)
	}

	resp, err := b.bedrockEmbeddingsRequest(toEmbedChunks)
	if err != nil {
		return err
	}

	for i, embedding := range resp.Embeddings {
		b.db.Add(allChunks[i], embedding)
	}
	err = b.db.SaveToFile()
	return err
}

func (b *Bedrock) bedrockEmbeddingsRequest(toEmbedChunks []string) (Response, error) {
	var resp Response
	payload := Request{
		Texts:     toEmbedChunks,
		InputType: "search_document",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Response{}, err
	}

	output, err := b.brc.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		Body:        payloadBytes,
		ModelId:     aws.String(titanEmbeddingModelID),
		ContentType: aws.String("application/json"),
	})

	if err != nil {
		return Response{}, err
	}

	err = json.Unmarshal(output.Body, &resp)

	if err != nil {
		return resp, err
	}
	return resp, nil
}
