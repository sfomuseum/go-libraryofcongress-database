package main

import (
	"context"
	_ "database/sql"
	"flag"
	"github.com/aaronland/go-sqlite"
	"github.com/aaronland/go-sqlite/database"
	loc_database "github.com/sfomuseum/go-libraryofcongress-database"
	loc_sqlite "github.com/sfomuseum/go-libraryofcongress-database/sqlite"
	"github.com/sfomuseum/go-libraryofcongress-database/sqlite/tables"
	loc_timings "github.com/sfomuseum/go-libraryofcongress-database/timings"
	"log"
	"os"
	"time"
)

func main() {

	dsn := flag.String("dsn", "libraryofcongress.db", "The SQLite DSN for the database you want to create.")

	lcsh_data := flag.String("lcsh-data", "", "The path to your LCSH CSV data.")
	lcnaf_data := flag.String("lcnaf-data", "", "The path to your LCNAF CSV data.")

	flag.Parse()

	ctx := context.Background()

	sqlite_db, err := database.NewDB(ctx, *dsn)

	if err != nil {
		log.Fatalf("Failed to create new database, %v", err)
	}

	err = sqlite_db.LiveHardDieFast()

	if err != nil {
		log.Fatalf("Failed to enable live hard, die fast settings, %v", err)
	}

	search_table, err := tables.NewSearchTableWithDatabase(ctx, sqlite_db)

	if err != nil {
		log.Fatalf("Failed to create search table, %v", err)
	}

	tables := []sqlite.Table{
		search_table,
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
	monitor, err := loc_timings.NewMonitor(ctx, d)

	if err != nil {
		log.Fatalf("Failed to create timings monitor, %v", err)
	}

	monitor.Start(ctx, os.Stdout)
	defer monitor.Stop(ctx)
	
	err = loc_sqlite.Index(ctx, data_sources, sqlite_db, tables, monitor)

	if err != nil {
		log.Fatalf("Failed to index sources, %v", err)
	}
}
