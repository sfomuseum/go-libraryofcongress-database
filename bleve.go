package database

import (
	"context"
	"fmt"
	"github.com/aaronland/go-pagination"
	"github.com/blevesearch/bleve"
	"log"
	"net/url"
	"os"
)

type BleveDatabase struct {
	LibraryOfCongressDatabase
	index bleve.Index
}

func init() {
	ctx := context.Background()
	RegisterLibraryOfCongressDatabase(ctx, "bleve", NewBleveDatabase)
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

func NewBleveDatabase(ctx context.Context, uri string) (LibraryOfCongressDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	path := u.Path

	index, err := NewBleveIndex(ctx, path)

	if err != nil {
		log.Fatalf("Failed to load Bleve index, %w", err)
	}

	bleve_db := &BleveDatabase{
		index: index,
	}

	return bleve_db, nil
}

func (bleve_db *BleveDatabase) Query(ctx context.Context, q string, pg_opts pagination.PaginationOptions) ([]*QueryResult, pagination.Pagination, error) {

	size := int(pg_opts.PerPage())
	from := int(pg_opts.PerPage() * pg_opts.PerPage())

	query := bleve.NewQueryStringQuery(q)
	req := bleve.NewSearchRequestOptions(query, size, from, false)

	rsp, err := bleve_db.index.Search(req)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to perform query, %w", err)
	}

	// https://pkg.go.dev/github.com/blevesearch/bleve#SearchResult

	log.Println(rsp.Total)
	log.Println(len(rsp.Hits))
	for _, d := range rsp.Hits {
		log.Println(d)
	}

	return nil, nil, fmt.Errorf("Not implemented")
}
