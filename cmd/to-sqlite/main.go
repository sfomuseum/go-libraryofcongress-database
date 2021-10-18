package main

import (
	"compress/bzip2"
	"context"
	_ "database/sql"
	"flag"
	"fmt"
	"github.com/aaronland/go-sqlite"
	"github.com/aaronland/go-sqlite/database"
	"github.com/sfomuseum/go-csvdict"
	"github.com/sfomuseum/go-libraryofcongress-database/sqlite/tables"
	"io"
	"log"
	"os"
	"path/filepath"
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

	data_sources := make(map[string]io.Reader)
	data_paths := make(map[string]string)

	if *lcsh_data != "" {
		data_paths["lcsh"] =  *lcsh_data
	}
	
	if *lcnaf_data != "" {		
		data_paths["lcnaf"] = *lcnaf_data
	}

	for source, path := range data_paths {

		r, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s, %v", path, err)
		}

		defer r.Close()

		ext := filepath.Ext(path)

		switch ext {
		case ".bz2":
			 data_sources[source] = bzip2.NewReader(r)
		default:
			data_sources[source] = r
		}		
	}

	//

	for source, r := range data_sources {

		err := index(ctx, source, sqlite_db, tables, r)

		if err != nil {
			log.Fatalf("Failed to index %s, %v", source, err)
		}

		log.Printf("Finished indexing %s\n", data_paths[source])
	}

}

func index(ctx context.Context, source string, db sqlite.Database, tables []sqlite.Table, r io.Reader) error {

	csv_r, err := csvdict.NewReader(r)

	if err != nil {
		return fmt.Errorf("Failed to create CSV reader for %s, %w", source, err)
	}

	for {
		row, err := csv_r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		row["source"] = source

		for _, t := range tables {

			err := t.IndexRecord(ctx, db, row)

			if err != nil {
				return fmt.Errorf("Failed to index %v in table %s, %w", row, t.Name(), err)
			}
		}

	}

	return nil
}
