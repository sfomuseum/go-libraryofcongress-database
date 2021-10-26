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
		
		log.Println("INDEX ", row)
		go monitor.Signal(ctx)

		return nil
	}

	return src.Index(ctx, cb)
}
