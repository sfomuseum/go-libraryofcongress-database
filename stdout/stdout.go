package stdout

import (
	"context"
	"fmt"
	"log"

	"github.com/aaronland/go-pagination"
	"github.com/sfomuseum/go-libraryofcongress-database"
	"github.com/sfomuseum/go-timings"
)

type StdoutDatabase struct {
	database.LibraryOfCongressDatabase
}

func init() {
	ctx := context.Background()
	database.RegisterLibraryOfCongressDatabase(ctx, "stdout", NewStdoutDatabase)
}

func NewStdoutDatabase(ctx context.Context, uri string) (database.LibraryOfCongressDatabase, error) {

	stdout_db := &StdoutDatabase{}

	return stdout_db, nil
}

func (stdout_db *StdoutDatabase) Index(ctx context.Context, sources []*database.Source, monitor timings.Monitor) error {

	for _, src := range sources {

		err := stdout_db.indexSource(ctx, src, monitor)

		if err != nil {
			return fmt.Errorf("Failed to index %s, %w", src.Label, err)
		}
	}

	return nil
}

func (stdout_db *StdoutDatabase) indexSource(ctx context.Context, src *database.Source, monitor timings.Monitor) error {

	cb := func(ctx context.Context, row map[string]string) error {

		log.Println(row)

		go monitor.Signal(ctx)
		return nil
	}

	return src.Index(ctx, cb)
}

func (stdout_db *StdoutDatabase) Query(ctx context.Context, q string, pg_opts pagination.Options) ([]*database.QueryResult, pagination.Results, error) {

	return nil, nil, fmt.Errorf("Not implemented")
}
