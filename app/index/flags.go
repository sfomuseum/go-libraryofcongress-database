package index

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var database_uri string
var lcsh_data string
var lcnaf_data string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("loc")

	fs.StringVar(&database_uri, "database-uri", "", "A valid sfomuseum/go-libraryofcongress-database URI.")

	fs.StringVar(&lcsh_data, "lcsh-data", "", "The path to your LCSH CSV data. If '-' then data will be read from STDIN.")
	fs.StringVar(&lcnaf_data, "lcnaf-data", "", "The path to your LCNAF CSV data. If '-' then data will be read from STDIN.")

	return fs
}
