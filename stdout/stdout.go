package stdout

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/aaronland/go-pagination"
	"github.com/sfomuseum/go-csvdict"
	"github.com/sfomuseum/go-libraryofcongress-database"
	"github.com/sfomuseum/go-timings"
)

type StdoutDatabase struct {
	database.LibraryOfCongressDatabase
	csv_writer *csvdict.Writer
	mu         *sync.RWMutex
}

func init() {
	ctx := context.Background()
	database.RegisterLibraryOfCongressDatabase(ctx, "stdout", NewStdoutDatabase)
}

func NewStdoutDatabase(ctx context.Context, uri string) (database.LibraryOfCongressDatabase, error) {

	mu := new(sync.RWMutex)

	stdout_db := &StdoutDatabase{
		mu: mu,
	}

	return stdout_db, nil
}

func (stdout_db *StdoutDatabase) Index(ctx context.Context, sources []*database.Source, monitor timings.Monitor) error {

	for _, src := range sources {

		err := stdout_db.indexSource(ctx, src, monitor)

		if err != nil {
			return fmt.Errorf("Failed to index %s, %w", src.Label, err)
		}
	}

	if stdout_db.csv_writer != nil {
		stdout_db.csv_writer.Flush()
	}

	return nil
}

func (stdout_db *StdoutDatabase) indexSource(ctx context.Context, src *database.Source, monitor timings.Monitor) error {

	cb := func(ctx context.Context, row map[string]string) error {

		stdout_db.mu.Lock()
		defer stdout_db.mu.Unlock()

		if stdout_db.csv_writer == nil {

			fieldnames := make([]string, 0)

			for k, _ := range row {
				fieldnames = append(fieldnames, k)
			}

			wr, err := csvdict.NewWriter(os.Stdout, fieldnames)

			if err != nil {
				return fmt.Errorf("Failed to create CSV writer, %w", err)
			}

			err = wr.WriteHeader()

			if err != nil {
				return fmt.Errorf("Failed to write CSV header, %w", err)
			}

			stdout_db.csv_writer = wr
		}

		err := stdout_db.csv_writer.WriteRow(row)

		if err != nil {
			return fmt.Errorf("Failed to write CSV row, %w", err)
		}

		go monitor.Signal(ctx)
		return nil
	}

	return src.Index(ctx, cb)
}

func (stdout_db *StdoutDatabase) Query(ctx context.Context, q string, pg_opts pagination.Options) ([]*database.QueryResult, pagination.Results, error) {

	return nil, nil, fmt.Errorf("Not implemented")
}
