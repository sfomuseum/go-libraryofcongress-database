package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	loc_database "github.com/sfomuseum/go-libraryofcongress-database"
	"github.com/sfomuseum/go-timings"
	"log"
)

func Index(ctx context.Context, data_sources []*loc_database.Source, bi esutil.BulkIndexer, monitor timings.Monitor) error {

	for _, src := range data_sources {

		err := index(ctx, src, bi, monitor)

		if err != nil {
			return fmt.Errorf("Failed to index %s, %v", src.Label, err)
		}

		log.Printf("Finished indexing %s\n", src.Label)
	}

	return nil
}

func index(ctx context.Context, src *loc_database.Source, bi esutil.BulkIndexer, monitor timings.Monitor) error {

	cb := func(ctx context.Context, row map[string]string) error {

		doc := map[string]string{
			"id":     row["id"],
			"label":  row["label"],
			"source": src.Label,
		}

		doc_id := row["id"]

		enc_doc, err := json.Marshal(doc)

		if err != nil {
			return fmt.Errorf("Failed to marshal %s, %v", doc_id, err)
		}

		// log.Println(string(enc_doc))
		// continue

		bulk_item := esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: doc_id,
			Body:       bytes.NewReader(enc_doc),

			OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
				// log.Printf("Indexed %s\n", path)
			},

			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
				if err != nil {
					log.Printf("ERROR: Failed to index %s, %s", doc_id, err)
				} else {
					log.Printf("ERROR: Failed to index %s, %s: %s", doc_id, res.Error.Type, res.Error.Reason)
				}
			},
		}

		err = bi.Add(ctx, bulk_item)

		if err != nil {
			log.Printf("Failed to schedule %s, %v", doc_id, err)
			return nil
		}

		go monitor.Signal(ctx)
		return nil
	}

	return src.Index(ctx, cb)
}
