package database

import (
	"compress/bzip2"
	"context"
	"fmt"
	"github.com/sfomuseum/go-csvdict"
	"io"
	"os"
	"path/filepath"
)

type SourceIndexCallback func(context.Context, map[string]string) error

type Source struct {
	Label  string
	Reader io.Reader
}

func (src *Source) Index(ctx context.Context, cb SourceIndexCallback) error {

	csv_r, err := csvdict.NewReader(src.Reader)

	if err != nil {
		return fmt.Errorf("Failed to create CSV reader for %s, %w", src.Label, err)
	}

	for {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		row, err := csv_r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		row["source"] = src.Label

		err = cb(ctx, row)

		if err != nil {
			return fmt.Errorf("Callback failed for %v, %w", row, err)
		}

	}

	return nil
}

func SourcesFromPaths(ctx context.Context, data_paths map[string]string) ([]*Source, error) {

	data_sources := make([]*Source, 0)

	for source, path := range data_paths {

		var r io.Reader

		fh, err := os.Open(path)

		if err != nil {
			return nil, fmt.Errorf("Failed to open %s, %v", path, err)
		}

		ext := filepath.Ext(path)

		switch ext {
		case ".bz2":
			r = bzip2.NewReader(fh)
		default:
			r = fh
		}

		src := &Source{
			Label:  source,
			Reader: r,
		}

		data_sources = append(data_sources, src)
	}

	return data_sources, nil
}
