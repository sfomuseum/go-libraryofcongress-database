package bleve

import (
	"context"
	"fmt"
	"github.com/blevesearch/bleve"
	"os"
)

type Doc struct {
	Id    string
	Label string
}

func NewBleveIndex(ctx context.Context, uri string) (bleve.Index, error) {

	var index bleve.Index

	_, err := os.Stat(uri)

	if err != nil {

		mapping := bleve.NewIndexMapping()

		locMapping := bleve.NewDocumentMapping()
		mapping.AddDocumentMapping("loc", locMapping)

		labelFieldMapping := bleve.NewTextFieldMapping()
		labelFieldMapping.Store = true

		locMapping.AddFieldMappingsAt("label", labelFieldMapping)

		i, err := bleve.New(uri, mapping)

		if err != nil {
			return nil, fmt.Errorf("Failed to create new index, %v", err)
		}

		index = i

	} else {

		i, err := bleve.Open(uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to open Bleve index, %w", err)
		}

		index = i
	}

	return index, nil
}
