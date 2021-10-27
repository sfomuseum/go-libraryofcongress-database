package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-pagination/countable"
	"github.com/sfomuseum/go-libraryofcongress-database"
	_ "github.com/sfomuseum/go-libraryofcongress-database/bleve"
	_ "github.com/sfomuseum/go-libraryofcongress-database/sql"
	"log"
	"strings"
)

func main() {

	db_uri := flag.String("database-uri", "", "A valid sfomuseum/go-libraryofcongress-database URI.")

	flag.Parse()

	ctx := context.Background()

	db, err := database.NewLibraryOfCongressDatabase(ctx, *db_uri)

	if err != nil {
		log.Fatalf("Failed to create database, %v", err)
	}

	pg_opts, err := countable.NewCountablePaginationOptions()

	if err != nil {
		log.Fatalf("Failed to create pagination options, %v", err)
	}

	cb := func(ctx context.Context, results []*database.QueryResult) error {

		for _, r := range results {
			fmt.Println(r)
		}

		return nil
	}

	q := strings.Join(flag.Args(), " ")

	err = database.QueryPaginated(ctx, db, q, pg_opts, cb)

	if err != nil {
		log.Fatalf("Failed to query for '%s', %v", q, err)
	}
}
