// index-es is a command-tool to index the `lcsh` and `lcnaf` data embedded in the `go-sfomuseum-libraryofcongress` package
// in an Elasticsearch index.
package main

import (
	"context"
	"flag"
	"github.com/cenkalti/backoff/v4"
	es "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	loc_database "github.com/sfomuseum/go-libraryofcongress-database"
	loc_elasticsearch "github.com/sfomuseum/go-libraryofcongress-database/elasticsearch"
	"github.com/sfomuseum/go-timings"
	"log"
	"os"
	"time"
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

	err = loc_elasticsearch.Index(ctx, data_sources, bi, monitor)

	if err != nil {
		log.Fatalf("Failed to index sources %v", err)
	}
	err = bi.Close(ctx)

	if err != nil {
		log.Fatalf("Failed to close indexer, %v", err)
	}

}
