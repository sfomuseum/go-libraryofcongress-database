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

func NewBleveDatabase(ctx context.Context, uri string) (LibraryOfCongressDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	path := u.Path

	_, err = os.Stat(path)

	var index bleve.Index

	if err != nil {

		mapping := bleve.NewIndexMapping()

		index, err = bleve.New(path, mapping)

		if err != nil {
			log.Fatalf("Failed to create Bleve index, %w", err)
		}

	} else {

		index, err = bleve.Open(path)

		if err != nil {
			log.Fatalf("Failed to open Bleve index, %w", err)
		}
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
