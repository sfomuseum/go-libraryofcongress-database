package main

import (
	"context"
	"log"

	"github.com/sfomuseum/go-libraryofcongress-database/app/query"
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
