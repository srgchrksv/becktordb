package becktordb

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
)

type Vector []float64

type VectorKey struct {
	Path      string
	TextChunk string
}

type VectorDB struct {
	Vectors      map[VectorKey]Vector
	DatabaseName string
	mu           sync.RWMutex
}

func NewVectorDB(databaseName string) *VectorDB {
	log.Printf("Init database: %s\n", databaseName)
	return &VectorDB{
		Vectors:      make(map[VectorKey]Vector),
		DatabaseName: databaseName,
		mu:           sync.RWMutex{},
	}
}

// Adds vector to database
func (db *VectorDB) Add(id VectorKey, vector Vector) {
	db.mu.Lock()
	defer db.mu.Unlock()

	log.Printf("Addint vector: %s\nText: %s", id.Path, id.TextChunk)
	db.Vectors[id] = vector
}

// Query for similarity
func (db *VectorDB) Query(query Vector, topK int) ([]VectorKey, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	scores := make(map[VectorKey]float64)
	err := db.LoadFromFile()
	if err != nil {
		return nil, err
	}
	for id, vector := range db.Vectors {
		scores[id] = CosineSimilarity(query, vector)
	}

	// Get top K results
	type kv struct {
		Key   VectorKey
		Value float64
	}
	var sorted []kv
	for k, v := range scores {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	var results []VectorKey
	for i := 0; i < topK && i < len(sorted); i++ {
		results = append(results, sorted[i].Key)
	}

	return results, err
}

// SaveToFile writes the VectorDB to a file
func (db *VectorDB) SaveToFile() error {
	file, err := os.Create(db.DatabaseName)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(db); err != nil {
		return fmt.Errorf("failed to encode data: %v", err)
	}

	return nil
}

// LoadFromFile loads the VectorDB from a file
func (db *VectorDB) LoadFromFile() error {
	file, err := os.Open(db.DatabaseName)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(db); err != nil {
		return fmt.Errorf("failed to decode data: %v", err)
	}

	return nil
}
