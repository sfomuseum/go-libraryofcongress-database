package api

import (
	"github.com/sfomuseum/go-libraryofcongress-database"
)

type Options struct {
	Database database.LibraryOfCongressDatabase
	PerPage  int64
}
