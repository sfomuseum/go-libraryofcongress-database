// index-es is a command-tool to index the `lcsh` and `lcnaf` data embedded in the `go-sfomuseum-libraryofcongress` package
// in an Elasticsearch index.
package main

import (
	"bytes"
	"compress/bzip2"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	es "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/sfomuseum/go-csvdict"
	"io"
	"log"
	"time"
	"os"
	"path/filepath"
)

func main() {

	es_endpoint := flag.String("elasticsearch-endpoint", "http://localhost:9200", "The Elasticsearch endpoint where data should be indexed.")
	es_index := flag.String("elasticsearch-index", "libraryofcongress", "The Elasticsearch index where data should be stored.")

	lcsh_data := flag.String("lcsh-data", "", "The path to your LCSH CSV data.")
	lcnaf_data := flag.String("lcnaf-data", "", "The path to your LCNAF CSV data.")	
	
	workers := flag.Int("workers", 10, "The number of concurrent workers to use when indexing data.")

	flag.Parse()

	ctx := context.Background()

	retry := backoff.NewExponentialBackOff()

	es_cfg := es.Config{
		Addresses: []string{*es_endpoint},

		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retry.Reset()
			}
			return retry.NextBackOff()
		},
		MaxRetries: 5,
	}

	es_client, err := es.NewClient(es_cfg)

	if err != nil {
		log.Fatalf("Failed to create ES client, %v", err)
	}

	_, err = es_client.Indices.Create(*es_index)

	if err != nil {
		log.Fatalf("Failed to create index, %v", err)
	}

	// https://github.com/elastic/go-elasticsearch/blob/master/_examples/bulk/indexer.go

	bi_cfg := esutil.BulkIndexerConfig{
		Index:         *es_index,
		Client:        es_client,
		NumWorkers:    *workers,
		FlushInterval: 30 * time.Second,
	}

	bi, err := esutil.NewBulkIndexer(bi_cfg)

	if err != nil {
		log.Fatalf("Failed to create bulk indexer, %v", err)
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

		err := index(ctx, bi, source, r)

		if err != nil {
			log.Fatalf("Failed to index %s, %v", source, err)
		}

		log.Printf("Finished indexing %s\n", data_paths[source])
	}

	err = bi.Close(ctx)

	if err != nil {
		log.Fatalf("Failed to close indexer, %v", err)
	}

}

func index(ctx context.Context, bi esutil.BulkIndexer, source string, r io.Reader) error {

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

		doc := map[string]string{
			"id":     row["id"],
			"label":  row["label"],
			"source": source,
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
			continue
		}
	}

	return nil
}
