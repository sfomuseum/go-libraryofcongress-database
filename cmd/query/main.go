package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sfomuseum/go-libraryofcongress-database/app/query"
	_ "github.com/sfomuseum/go-libraryofcongress-database/bleve"
	_ "github.com/sfomuseum/go-libraryofcongress-database/docstore"
	_ "github.com/sfomuseum/go-libraryofcongress-database/elasticsearch"
	_ "github.com/sfomuseum/go-libraryofcongress-database/sql"
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	err := query.Run(ctx, logger)

	if err != nil {
		logger.Fatalf("Failed to run query, %v", err)
	}
}
