package main

import (
	"context"
	"flag"
	"github.com/aaronland/go-pagination/countable"
	"github.com/sfomuseum/go-libraryofcongress-database"
	_ "github.com/sfomuseum/go-libraryofcongress-database/bleve"
	_ "github.com/sfomuseum/go-libraryofcongress-database/sql"
	"log"
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

	for _, q := range flag.Args() {
		db.Query(ctx, q, pg_opts)
	}

}
