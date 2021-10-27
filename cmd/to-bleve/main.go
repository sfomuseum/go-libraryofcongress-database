package main

import (
	"context"
	_ "database/sql"
	"flag"
	loc_database "github.com/sfomuseum/go-libraryofcongress-database"
	loc_bleve "github.com/sfomuseum/go-libraryofcongress-database/bleve"
	"github.com/sfomuseum/go-timings"
	"log"
	"os"
	"time"
)

func main() {

	path_index := flag.String("index", "libraryofcongress.db", "The path to the Bleve index you want to create.")

	lcsh_data := flag.String("lcsh-data", "", "The path to your LCSH CSV data.")
	lcnaf_data := flag.String("lcnaf-data", "", "The path to your LCNAF CSV data.")

	flag.Parse()

	ctx := context.Background()

	//

	index, err := loc_bleve.NewBleveIndex(ctx, *path_index)

	if err != nil {
		log.Fatalf("Failed to load Bleve index, %w", err)
	}

	//

	data_paths := make(map[string]string)

	if *lcsh_data != "" {
		data_paths["lcsh"] = *lcsh_data
	}

	if *lcnaf_data != "" {
		data_paths["lcnaf"] = *lcnaf_data
	}

	data_sources, err := loc_database.SourcesFromPaths(ctx, data_paths)

	if err != nil {
		log.Fatalf("Failed to derive database sources from paths, %v", err)
	}

	d := time.Second * 60
	monitor, err := timings.NewCounterMonitor(ctx, d)

	if err != nil {
		log.Fatalf("Failed to create timings monitor, %v", err)
	}

	monitor.Start(ctx, os.Stdout)
	defer monitor.Stop(ctx)

	err = loc_bleve.Index(ctx, data_sources, index, monitor)

	if err != nil {
		log.Fatalf("Failed to index sources, %v", err)
	}
}
