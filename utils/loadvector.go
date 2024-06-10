package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// LoadVectorDB loads a vector database from a file and returns a slice of float32 slices.
func LoadVectorDB(path string) [][]float32 {
	embeddings, err := ReadEmbeddings(path)
	if err != nil {
		fmt.Println("Error:", err)
		return [][]float32{}
	}

	// for i, embedding := range embeddings {
	// 	fmt.Printf("Embedding %d: %v\n", i, embedding)
	// }
	return embeddings
}

// ReadEmbeddings reads embeddings from a file and returns a slice of float32 slices.
func ReadEmbeddings(filename string) ([][]float32, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var embeddings [][]float32
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		embedding, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line: %v", err)
		}
		embeddings = append(embeddings, embedding)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return embeddings, nil
}

// parseLine parses a line of text into a slice of float32.
func parseLine(line string) ([]float32, error) {
	line = strings.Trim(line, "[]")
	parts := strings.Fields(line)
	var embedding []float32

	for _, part := range parts {
		value, err := strconv.ParseFloat(part, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse float: %v", err)
		}
		embedding = append(embedding, float32(value))
	}

	return embedding, nil
}
