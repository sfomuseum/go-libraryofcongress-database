package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/sfomuseum/go-csvdict"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-libraryofcongress-database"
)

func main() {

	var lcsh_data string
	var lcnaf_data string

	fs := flagset.NewFlagSet("loc")

	fs.StringVar(&lcsh_data, "lcsh-data", "", "The path to your LCSH CSV data.")
	fs.StringVar(&lcnaf_data, "lcnaf-data", "", "The path to your LCNAF CSV data.")

	flagset.Parse(fs)

	ctx := context.Background()

	data_paths := make(map[string]string)

	if lcsh_data != "" {
		data_paths["lcsh"] = lcsh_data
	}

	if lcnaf_data != "" {
		data_paths["lcnaf"] = lcnaf_data
	}

	data_sources, err := database.SourcesFromPaths(ctx, data_paths)

	if err != nil {
		log.Fatalf("Failed to derive database sources from paths, %v", err)
	}

	var csv_wr *csvdict.Writer
	wr := os.Stdout

	for _, s := range data_sources {

		r := s.Reader
		defer r.Close()

		csv_r, err := csvdict.NewReader(r)

		if err != nil {
			log.Fatalf("Failed to create CSV reader for %s, %v", s.Label, err)
		}

		for {
			row, err := csv_r.Read()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Failed to read row, %v", err)
			}

			out := map[string]string{
				"Id":     row["id"],
				"Label":  row["label"],
				"Source": s.Label,
			}

			if csv_wr == nil {

				fieldnames := make([]string, 0)

				for k, _ := range out {
					fieldnames = append(fieldnames, k)
				}

				w, err := csvdict.NewWriter(wr, fieldnames)

				if err != nil {
					log.Fatalf("Failed to create new CSV writer, %v", err)
				}

				csv_wr = w

				err = csv_wr.WriteHeader()

				if err != nil {
					log.Fatalf("Failed to write CSV header, %v", err)
				}
			}

			err = csv_wr.WriteRow(out)

			if err != nil {
				log.Fatalf("Failed to write row, %v", err)
			}

		}
	}

	csv_wr.Flush()
}
