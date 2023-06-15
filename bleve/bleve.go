package bleve

import (
	"context"
	"fmt"
	"os"

	"github.com/blevesearch/bleve"	
)

type Doc struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Source string `json:"source"`
}

func (d *Doc) String() string {
	return fmt.Sprintf("%s:%s %s", d.Source, d.Id, d.Label)
}

func NewBleveIndex(ctx context.Context, uri string) (bleve.Index, error) {

	var index bleve.Index

	_, err := os.Stat(uri)

	if err != nil {

		mapping := bleve.NewIndexMapping()

		docMapping := bleve.NewDocumentMapping()
		mapping.AddDocumentMapping("doc", docMapping)

		labelFieldMapping := bleve.NewTextFieldMapping()
		labelFieldMapping.Store = true

		sourceFieldMapping := bleve.NewTextFieldMapping()
		sourceFieldMapping.Store = true

		docMapping.AddFieldMappingsAt("label", labelFieldMapping)
		docMapping.AddFieldMappingsAt("source", sourceFieldMapping)

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
