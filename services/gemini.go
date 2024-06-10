package services

import (
	"context"
	"os"
	"sync"

	"github.com/srgchrksv/becktordb/becktordb"
	"google.golang.org/api/option"

	"github.com/google/generative-ai-go/genai"
	"github.com/srgchrksv/becktordb/utils"
)

type Gemini struct {
	db *becktordb.VectorDB
	mu sync.RWMutex
}

func (g *Gemini) Init() *Gemini {
	return &Gemini{db: g.db, mu: sync.RWMutex{}}
}

// GeminiEmbeddingsQuery embed query and run vector similarity, return closest 2 chunks
func (g *Gemini) GeminiEmbeddingsQuery(query string, topK int) ([]becktordb.VectorKey, error) {
	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, err
	}
	defer client.Close()
	// For embeddings, use the embedding-001 model
	em := client.EmbeddingModel("embedding-001")
	res, err := em.EmbedContent(ctx, genai.Text(query))

	if err != nil {
		return nil, err
	}
	response, err := g.db.Query(utils.Float32ToFloat64(res.Embedding.Values), topK)
	return response, err
}

func (g *Gemini) GeminiEmbeddingsPath(directories []string, chunkSize int) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return err
	}
	defer client.Close()
	em := client.EmbeddingModel("embedding-001")
	b := em.NewBatch()

	allChunks, err := Chunking(directories, chunkSize)
	if err != nil {
		return err
	}
	for _, chunk := range allChunks {
		b = b.AddContent(genai.Text(chunk.TextChunk))
	}
	res, err := em.BatchEmbedContents(ctx, b)
	if err != nil {
		return err
	}

	for i, chunk := range res.Embeddings {
		g.db.Add(allChunks[i], utils.Float32ToFloat64(chunk.Values))
	}
	err = g.db.SaveToFile()

	return err
}
