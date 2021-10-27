package bleve

import (
	"context"
	"fmt"
	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/countable"
	"github.com/blevesearch/bleve"
	"github.com/sfomuseum/go-libraryofcongress-database"
	_ "log"
	"net/url"
)

type BleveDatabase struct {
	database.LibraryOfCongressDatabase
	index bleve.Index
}

func init() {
	ctx := context.Background()
	database.RegisterLibraryOfCongressDatabase(ctx, "bleve", NewBleveDatabase)
}

func NewBleveDatabase(ctx context.Context, uri string) (database.LibraryOfCongressDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	path := u.Path

	index, err := NewBleveIndex(ctx, path)

	if err != nil {
		return nil, fmt.Errorf("Failed to load Bleve index, %w", err)
	}

	bleve_db := &BleveDatabase{
		index: index,
	}

	return bleve_db, nil
}

func (bleve_db *BleveDatabase) Query(ctx context.Context, q string, pg_opts pagination.PaginationOptions) ([]*database.QueryResult, pagination.Pagination, error) {

	page := pg_opts.Page()
	per_page := pg_opts.PerPage()

	size := int(per_page)
	from := int(per_page * (page - 1))

	query := bleve.NewQueryStringQuery(q)
	req := bleve.NewSearchRequestOptions(query, size, from, false)

	req.Fields = []string{
		"label",
		"source",
	}

	rsp, err := bleve_db.index.Search(req)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to perform query, %w", err)
	}

	results := make([]*database.QueryResult, 0)

	for _, d := range rsp.Hits {

		id := d.ID
		fields := d.Fields

		r := &database.QueryResult{
			Id:     id,
			Label:  fields["label"].(string),
			Source: fields["source"].(string),
		}

		results = append(results, r)
	}

	total := int64(rsp.Total)

	pg, err := countable.NewPaginationFromCountWithOptions(pg_opts, total)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to derive pagination, %w", err)
	}

	return results, pg, nil
}
