package docstore

import (
	"context"
	"fmt"
	"io"
	"regexp"

	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/cursor"
	aa_docstore "github.com/aaronland/gocloud-docstore"
	"github.com/sfomuseum/go-libraryofcongress-database"
	"github.com/sfomuseum/go-timings"
	gc_docstore "gocloud.dev/docstore"
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
	re_lcnaf   *regexp.Regexp
	re_lcsh    *regexp.Regexp
}

func NewDocstoreDatabase(ctx context.Context, uri string) (database.LibraryOfCongressDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create collection, %w", err)
	}

	re_lcnaf, err := regexp.Compile(`^n\d+$`)

	if err != nil {
		return nil, fmt.Errorf("Failed to compile LCNAF pattern, %w", err)
	}

	re_lcsh, err := regexp.Compile(`^sh\d+$`)

	if err != nil {
		return nil, fmt.Errorf("Failed to compile LCSH pattern, %w", err)
	}

	db := &DocstoreDatabase{
		collection: col,
		re_lcnaf:   re_lcnaf,
		re_lcsh:    re_lcsh,
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

	results := make([]*database.QueryResult, 0)

	limit := int(pg_opts.PerPage())

	q := db.collection.Query()

	if db.re_lcnaf.MatchString(query) || db.re_lcsh.MatchString(query) {
		q = q.Where("Id", "=", query)
	} else {
		q = q.Where("Label", "=", query)
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
			Id:     doc.Id,
			Label:  doc.Label,
			Source: doc.Source,
		}

		results = append(results, qr)
	}

	if len(results) > 0 {
		next_cursor = results[len(results)-1].Id
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
