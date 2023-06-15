package query

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var database_uri string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("loc")

	fs.StringVar(&database_uri, "database-uri", "sql://sqlite3?dsn=data/lcsh.db", "A valid sfomuseum/go-libraryofcongress-database URI.")
	return fs
}
