package docstore

import (
	"context"
	"fmt"
	"io"
	
	"github.com/sfomuseum/go-libraryofcongress-database"
	aa_docstore "github.com/aaronland/gocloud-docstore"
	gc_docstore "gocloud.dev/docstore"
	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/cursor"	
	"github.com/sfomuseum/go-timings"
	
)

func init() {

	ctx := context.Background()

	database.RegisterLibraryOfCongressDatabase(ctx, "awsdynamodb", NewDocstoreDatabase)

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		database.RegisterLibraryOfCongressDatabase(ctx, scheme, NewDocstoreDatabase)
	}
}

type DocstoreDatabase struct {
	database.LibraryOfCongressDatabase
	collection *gc_docstore.Collection
}

func NewDocstoreDatabase(ctx context.Context, uri string) (database.LibraryOfCongressDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create collection, %w", err)
	}

	db := &DocstoreDatabase{
		collection: col,
	}

	return db, nil
}

func (db *DocstoreDatabase) Index(ctx context.Context, sources []*database.Source, monitor timings.Monitor) error {

	for _, src := range sources {

		err := db.indexSource(ctx, src, monitor)

		if err != nil {
			return fmt.Errorf("Failed to index %s, %v", src.Label, err)
		}
	}

	return nil
}

func (db *DocstoreDatabase) Query(ctx context.Context, query string, pg_opts pagination.Options) ([]*database.QueryResult, pagination.Results, error) {

	var previous_cursor string
	var next_cursor string

	previous_cursor = pg_opts.Pointer().(string)
	
	results := make([]*database.QueryResult, 0)

	limit := int(pg_opts.PerPage())
	
	q := db.collection.Query()
	q = q.Where("Label", "=", query)

	if previous_cursor != "" {
		q = q.Where("Id", ">", previous_cursor)
	}
	
	q = q.Limit(limit)
	
	iter := q.Get(ctx)
	defer iter.Stop()

	for {
		
		var doc Document
		err := iter.Next(ctx, &doc)
		
		if err != nil {
			
			if err == io.EOF {
				break
			}

			return nil, nil, fmt.Errorf("Failed to retrieve next item in iterator, %w", err)
		}

		qr := &database.QueryResult{
			Id: doc.Id,
			Label: doc.Label,
			Source: doc.Source,
		}

		results = append(results, qr)
	}

	if len(results) > 0 {
		next_cursor = results[ len(results) - 1 ].Id
	}
	
	pg_results, err := cursor.NewPaginationFromCursors(previous_cursor, next_cursor)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to build pagination results, %w", err)
	}
	
	return results, pg_results, nil
}

func (db *DocstoreDatabase) indexSource(ctx context.Context, src *database.Source, monitor timings.Monitor) error {

	cb := func(ctx context.Context, row map[string]string) error {

		doc := &Document{
			Id:     row["id"],
			Label:  row["label"],
			Source: src.Label,
		}

		err := db.collection.Put(ctx, doc)

		if err != nil {
			return fmt.Errorf("Failed to index row, %w", err)
		}

		// log.Println(doc)
		go monitor.Signal(ctx)

		return nil
	}

	return src.Index(ctx, cb)
}
