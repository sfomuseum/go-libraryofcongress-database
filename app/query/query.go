package query

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/aaronland/go-pagination/countable"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-libraryofcongress-database"
)

func Run(ctx context.Context, logger *log.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *log.Logger) error {

	flagset.Parse(fs)

	db, err := database.NewLibraryOfCongressDatabase(ctx, database_uri)

	if err != nil {
		return fmt.Errorf("Failed to create database, %w", err)
	}

	pg_opts, err := countable.NewCountableOptions()

	if err != nil {
		return fmt.Errorf("Failed to create pagination options, %v", err)
	}

	cb := func(ctx context.Context, results []*database.QueryResult) error {

		for _, r := range results {
			fmt.Println(r)
		}

		return nil
	}

	q := strings.Join(fs.Args(), " ")

	err = database.QueryPaginated(ctx, db, q, pg_opts, cb)

	if err != nil {
		return fmt.Errorf("Failed to query for '%s', %v", q, err)
	}

	return nil
}
