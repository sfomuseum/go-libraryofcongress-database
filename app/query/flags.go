package query

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var database_uri string
var cursor_pagination bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("loc")

	fs.StringVar(&database_uri, "database-uri", "", "A valid sfomuseum/go-libraryofcongress-database URI.")
	fs.BoolVar(&cursor_pagination, "cursor-pagination", false, "Signal that pagination is cursor-based rather than countable.")
	return fs
}
