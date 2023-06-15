package sqlite

import (
	"context"
	"fmt"
	"log"

	"github.com/aaronland/go-sqlite"
	loc_database "github.com/sfomuseum/go-libraryofcongress-database"
	"github.com/sfomuseum/go-timings"	
)

func Index(ctx context.Context, sources []*loc_database.Source, sqlite_db sqlite.Database, tables []sqlite.Table, monitor timings.Monitor) error {

	for _, src := range sources {

		err := index(ctx, src, sqlite_db, tables, monitor)

		if err != nil {
			return fmt.Errorf("Failed to index %s, %v", src.Label, err)
		}

		log.Printf("Finished indexing %s\n", src.Label)
	}

	return nil
}

func index(ctx context.Context, src *loc_database.Source, db sqlite.Database, tables []sqlite.Table, monitor timings.Monitor) error {

	cb := func(ctx context.Context, row map[string]string) error {

		for _, t := range tables {

			err := t.IndexRecord(ctx, db, row)

			if err != nil {
				return fmt.Errorf("Failed to index %v in table %s, %w", row, t.Name(), err)
			}

			go monitor.Signal(ctx)
		}

		return nil
	}

	return src.Index(ctx, cb)
}
