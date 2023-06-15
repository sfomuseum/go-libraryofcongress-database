package database

import (
	"context"
	"fmt"

	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/cursor"
)

type QueryResult struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Source string `json:"source"`
}

func (r *QueryResult) String() string {
	return fmt.Sprintf("%s:%s %s", r.Source, r.Id, r.Label)
}

type QueryPaginatedCallbackFunc func(context.Context, []*QueryResult) error

func QueryPaginated(ctx context.Context, db LibraryOfCongressDatabase, q string, pg_opts pagination.Options, cb QueryPaginatedCallbackFunc) error {

	cursor_pagination := false

	switch pg_opts.(type) {
	case *cursor.CursorOptions:
		cursor_pagination = true
	default:
		//
	}

	pages := int64(-1)

	for {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		rsp, pg, err := db.Query(ctx, q, pg_opts)

		if err != nil {
			return err
		}

		err = cb(ctx, rsp)

		if err != nil {
			return err
		}

		if cursor_pagination {

			break

		} else {

			if pages == -1 {
				pages = pg.Pages()
			}

			page := pg_opts.Pointer().(int64)
			page += 1

			if page <= pages {
				pg_opts.Pointer(page)
			} else {
				break
			}
		}
	}

	return nil

}
