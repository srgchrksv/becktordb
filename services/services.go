package services

import (
	"fmt"

	"github.com/srgchrksv/becktordb/becktordb"
)

type BectorDBservice struct {
	db      *becktordb.VectorDB
	Gemini  Gemini
	Bedrock Bedrock
}

func NewBectorDBservice(db *becktordb.VectorDB) *BectorDBservice {
	fmt.Println("creating new service")
	return &BectorDBservice{
		db:      db,
		Gemini:  Gemini{db: db},
		Bedrock: Bedrock{db: db},
	}
}
