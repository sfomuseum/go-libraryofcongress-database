package main

import (
	_ "github.com/sfomuseum/go-csvdict"
	"github.com/sfomuseum/go-flags/flagset"
)

func main() {

	var lcsh_data string
	var lcnaf_data string

	fs := flagset.NewFlagSet("loc")

	fs.StringVar(&lcsh_data, "lcsh-data", "", "The path to your LCSH CSV data.")
	fs.StringVar(&lcnaf_data, "lcnaf-data", "", "The path to your LCNAF CSV data.")

	flagset.Parse(fs)

}
