package index

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-libraryofcongress-database"
	"github.com/sfomuseum/go-timings"
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

	data_paths := make(map[string]string)

	if lcsh_data != "" {
		data_paths["lcsh"] = lcsh_data
	}

	if lcnaf_data != "" {
		data_paths["lcnaf"] = lcnaf_data
	}

	data_sources, err := database.SourcesFromPaths(ctx, data_paths)

	if err != nil {
		return fmt.Errorf("Failed to derive database sources from paths, %v", err)
	}

	monitor, err := timings.NewMonitor(ctx, "counter://PT60S")

	if err != nil {
		return fmt.Errorf("Failed to create timings monitor, %v", err)
	}

	monitor.Start(ctx, os.Stdout)
	defer monitor.Stop(ctx)

	err = db.Index(ctx, data_sources, monitor)

	if err != nil {
		return fmt.Errorf("Failed to index sources, %v", err)
	}

	return nil
}
