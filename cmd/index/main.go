package main

import (
	"context"
	"log"

	"github.com/sfomuseum/go-libraryofcongress-database/app/index"
	_ "github.com/sfomuseum/go-libraryofcongress-database/sql"
	_ "github.com/sfomuseum/go-libraryofcongress-database/stdout"
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	err := index.Run(ctx, logger)

	if err != nil {
		logger.Fatalf("Failed to run indexer, %v", err)
	}
}
