// The `to-bleve` tool will index CSV data produced by the tools in `sfomuseum/go-libraryofcongress` in an Bleve index.
package main

import (
	"context"
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sfomuseum/go-libraryofcongress-database"
	_ "github.com/sfomuseum/go-libraryofcongress-database/bleve"
	_ "github.com/sfomuseum/go-libraryofcongress-database/elasticsearch"
	_ "github.com/sfomuseum/go-libraryofcongress-database/sql"
	"github.com/sfomuseum/go-timings"
	"log"
	"os"
)

func main() {

	database_uri := flag.String("database-uri", "", "...")

	lcsh_data := flag.String("lcsh-data", "", "The path to your LCSH CSV data.")
	lcnaf_data := flag.String("lcnaf-data", "", "The path to your LCNAF CSV data.")

	flag.Parse()

	ctx := context.Background()

	db, err := database.NewLibraryOfCongressDatabase(ctx, *database_uri)

	if err != nil {
		log.Fatalf("Failed to create database, %w", err)
	}

	data_paths := make(map[string]string)

	if *lcsh_data != "" {
		data_paths["lcsh"] = *lcsh_data
	}

	if *lcnaf_data != "" {
		data_paths["lcnaf"] = *lcnaf_data
	}

	data_sources, err := database.SourcesFromPaths(ctx, data_paths)

	if err != nil {
		log.Fatalf("Failed to derive database sources from paths, %v", err)
	}

	monitor, err := timings.NewMonitor(ctx, "counter://PT60S")

	if err != nil {
		log.Fatalf("Failed to create timings monitor, %v", err)
	}

	monitor.Start(ctx, os.Stdout)
	defer monitor.Stop(ctx)

	err = db.Index(ctx, data_sources, monitor)

	if err != nil {
		log.Fatalf("Failed to index sources, %v", err)
	}
}
