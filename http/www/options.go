package www

import (
	"github.com/sfomuseum/go-libraryofcongress-database"
	"html/template"
)

type Options struct {
	Database  database.LibraryOfCongressDatabase
	Templates *template.Template
	PerPage   int64
}
