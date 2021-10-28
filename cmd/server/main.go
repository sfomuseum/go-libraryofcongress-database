package main

// To do: Move all of this code in to a generic 'app' package
// so that it can be used with custom database drivers

import (
	"context"
	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-server"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-libraryofcongress-database"
	_ "github.com/sfomuseum/go-libraryofcongress-database/bleve"
	"github.com/sfomuseum/go-libraryofcongress-database/http/api"
	"github.com/sfomuseum/go-libraryofcongress-database/http/www"
	_ "github.com/sfomuseum/go-libraryofcongress-database/sql"
	"github.com/sfomuseum/go-libraryofcongress-database/templates/html"
	"html/template"
	"log"
	"net/http"
)

func main() {

	fs := flagset.NewFlagSet("libraryofcongress")

	db_uri := fs.String("database-uri", "sql://sqlite3?dsn=data/lcsh.db", "A valid sfomuseum/go-libraryofcongress-database URI.")
	server_uri := fs.String("server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI.")

	per_page := fs.Int64("per-page", 20, "The number of results to return per page")

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "LIBRARYOFCONGRESS")

	if err != nil {
		log.Fatalf("Failed to assign flags from environment variables, %v", err)
	}

	ctx := context.Background()

	db, err := database.NewLibraryOfCongressDatabase(ctx, *db_uri)

	if err != nil {
		log.Fatalf("Failed to create database, %v", err)
	}

	t := template.New("libraryofcongress")
	t, err = t.ParseFS(html.FS, "*.html")

	if err != nil {
		log.Fatalf("Failed to load templates, %v", err)
	}

	mux := http.NewServeMux()

	err = bootstrap.AppendAssetHandlers(mux)

	if err != nil {
		log.Fatalf("Failed to append Bootstrap assets handler, %v", err)
	}

	bootstrap_opts := bootstrap.DefaultBootstrapOptions()

	www_opts := &www.Options{
		Database:  db,
		Templates: t,
		PerPage:   *per_page,
	}

	query_handler, err := www.QueryHandler(www_opts)

	if err != nil {
		log.Fatalf("Failed to create query handler, %v", err)
	}

	query_handler = bootstrap.AppendResourcesHandler(query_handler, bootstrap_opts)

	mux.Handle("/", query_handler)

	api_opts := &api.Options{
		Database: db,
		PerPage:  *per_page,
	}

	api_query_handler, err := api.QueryHandler(api_opts)

	if err != nil {
		log.Fatalf("Failed to create API query handler, %v", err)
	}

	mux.Handle("/api/query", api_query_handler)

	s, err := server.NewServer(ctx, *server_uri)

	if err != nil {
		log.Fatalf("Failed to create new server, %v", err)
	}

	log.Printf("Listening on %s\n", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		log.Fatalf("Failed to listen and serve, %v", err)
	}
}
