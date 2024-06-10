package services

import (
	"fmt"
	"io/ioutil"

	"github.com/srgchrksv/becktordb/becktordb"

	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Chunk struct {
	Path  string
	Chunk string
}

func Chunking(directories []string, chunkSize int) ([]becktordb.VectorKey, error) {
	var all_chunks []becktordb.VectorKey

	docPaths, err := GetAllFiles(directories)
	if err != nil {
		fmt.Println("Error retrieving files:", err)
		return nil, err

	}
	for _, path := range docPaths {
		content, err := LoadDocument(path)
		if err != nil {
			fmt.Println("Error loading document:", err)
			continue
		}
		chunks := chunkText(content, chunkSize)
		for _, chunk := range chunks {
			all_chunks = append(all_chunks, becktordb.VectorKey{Path: path, TextChunk: chunk})
		}
	}

	return all_chunks, err
}

// GetAllFiles retrieves all file paths from the specified directories
func GetAllFiles(paths []string) ([]string, error) {
	var files []string
	for _, path := range paths {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}

// LoadDocument loads a document from the given path and returns its content as a string
func LoadDocument(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// chunkText splits the text into chunks of at least minWords, ensuring each chunk ends with a full sentence.
func chunkText(text string, minWords int) []string {
	// Regular expression to split the text into sentences
	re := regexp.MustCompile(`[.!?](?:\s+|$)`)
	sentences := re.Split(text, -1)

	var chunks []string
	var currentChunk []string
	var currentWordCount int

	for _, sentence := range sentences {
		sentenceWordCount := len(strings.Fields(sentence))

		if currentWordCount+sentenceWordCount >= minWords {
			chunks = append(chunks, strings.Join(currentChunk, " "))
			currentChunk = []string{sentence}
			currentWordCount = sentenceWordCount
		} else {
			currentChunk = append(currentChunk, sentence)
			currentWordCount += sentenceWordCount
		}
	}

	// Add the last chunk
	if len(currentChunk) > 0 {
		chunks = append(chunks, strings.Join(currentChunk, " "))
	}

	return chunks
}
