package bleve

import (
	"context"
	"fmt"
	"github.com/blevesearch/bleve"
	loc_database "github.com/sfomuseum/go-libraryofcongress-database"
	"github.com/sfomuseum/go-timings"
	"log"
)

func Index(ctx context.Context, sources []*loc_database.Source, bl_index bleve.Index, monitor timings.Monitor) error {

	for _, src := range sources {

		err := index(ctx, src, bl_index, monitor)

		if err != nil {
			return fmt.Errorf("Failed to index %s, %v", src.Label, err)
		}

		log.Printf("Finished indexing %s\n", src.Label)
	}

	return nil
}

func index(ctx context.Context, src *loc_database.Source, bl_index bleve.Index, monitor timings.Monitor) error {

	cb := func(ctx context.Context, row map[string]string) error {

		// https://blevesearch.com/docs/Index-Mapping/

		doc := &Doc{
			Id:     row["id"],
			Label:  row["label"],
			Source: src.Label,
		}

		err := bl_index.Index(row["id"], doc)

		if err != nil {
			return fmt.Errorf("Failed to index row, %w", err)
		}

		// log.Println(doc)
		go monitor.Signal(ctx)

		return nil
	}

	return src.Index(ctx, cb)
}
