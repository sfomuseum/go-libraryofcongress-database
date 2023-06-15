// The `to-bleve` tool will index CSV data produced by the tools in `sfomuseum/go-libraryofcongress` in an Bleve index.
package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sfomuseum/go-libraryofcongress-database/app/index"
	_ "github.com/sfomuseum/go-libraryofcongress-database/bleve"
	_ "github.com/sfomuseum/go-libraryofcongress-database/docstore"
	_ "github.com/sfomuseum/go-libraryofcongress-database/elasticsearch"
	_ "github.com/sfomuseum/go-libraryofcongress-database/sql"
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	err := index.Run(ctx, logger)

	if err != nil {
		logger.Fatalf("Failed to run indexer, %v", err)
	}
}
