package server

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var database_uri string
var server_uri string
var per_page int64

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("loc")

	fs.StringVar(&database_uri, "database-uri", "", "A valid sfomuseum/go-libraryofcongress-database URI.")
	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI.")
	fs.Int64Var(&per_page, "per-page", 20, "The number of results to return per page")

	return fs
}
