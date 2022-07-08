package database

import (
	"context"
	"fmt"
	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-roster"
	"github.com/sfomuseum/go-timings"
	"net/url"
	"sort"
	"strings"
)

type LibraryOfCongressDatabase interface {
	Index(context.Context, []*Source, timings.Monitor) error
	Query(context.Context, string, pagination.Options) ([]*QueryResult, pagination.Results, error)
}

type LibraryOfCongressInitializeFunc func(ctx context.Context, uri string) (LibraryOfCongressDatabase, error)

var libraryofcongress_databases roster.Roster

func ensureLibraryofcongressRoster() error {

	if libraryofcongress_databases == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		libraryofcongress_databases = r
	}

	return nil
}

func RegisterLibraryOfCongressDatabase(ctx context.Context, scheme string, f LibraryOfCongressInitializeFunc) error {

	err := ensureLibraryofcongressRoster()

	if err != nil {
		return err
	}

	return libraryofcongress_databases.Register(ctx, scheme, f)
}

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureLibraryofcongressRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range libraryofcongress_databases.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func NewLibraryOfCongressDatabase(ctx context.Context, uri string) (LibraryOfCongressDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := libraryofcongress_databases.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(LibraryOfCongressInitializeFunc)
	return f(ctx, uri)
}
