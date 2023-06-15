package server

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-server"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-libraryofcongress-database"
	"github.com/sfomuseum/go-libraryofcongress-database/http/api"
	"github.com/sfomuseum/go-libraryofcongress-database/http/www"
	"github.com/sfomuseum/go-libraryofcongress-database/templates/html"
)

func Run(ctx context.Context, logger *log.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *log.Logger) error {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "LIBRARYOFCONGRESS")

	if err != nil {
		return fmt.Errorf("Failed to assign flags from environment variables, %w", err)
	}

	db, err := database.NewLibraryOfCongressDatabase(ctx, database_uri)

	if err != nil {
		return fmt.Errorf("Failed to create database, %w", err)
	}

	t := template.New("libraryofcongress")
	t, err = t.ParseFS(html.FS, "*.html")

	if err != nil {
		return fmt.Errorf("Failed to load templates, %w", err)
	}

	mux := http.NewServeMux()

	bootstrap_opts := bootstrap.DefaultBootstrapOptions()

	err = bootstrap.AppendAssetHandlers(mux, bootstrap_opts)

	if err != nil {
		return fmt.Errorf("Failed to append Bootstrap assets handler, %w", err)
	}

	www_opts := &www.Options{
		Database:  db,
		Templates: t,
		PerPage:   per_page,
	}

	query_handler, err := www.QueryHandler(www_opts)

	if err != nil {
		return fmt.Errorf("Failed to create query handler, %w", err)
	}

	query_handler = bootstrap.AppendResourcesHandler(query_handler, bootstrap_opts)

	mux.Handle("/", query_handler)

	api_opts := &api.Options{
		Database: db,
		PerPage:  per_page,
	}

	api_query_handler, err := api.QueryHandler(api_opts)

	if err != nil {
		return fmt.Errorf("Failed to create API query handler, %w", err)
	}

	mux.Handle("/api/query", api_query_handler)

	s, err := server.NewServer(ctx, server_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new server, %w", err)
	}

	logger.Printf("Listening on %s\n", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		return fmt.Errorf("Failed to listen and serve, %w", err)
	}

	return nil
}
