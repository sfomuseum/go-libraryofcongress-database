package www

import (
	"html/template"

	"github.com/sfomuseum/go-libraryofcongress-database"
)

type Options struct {
	Database  database.LibraryOfCongressDatabase
	Templates *template.Template
	PerPage   int64
}
