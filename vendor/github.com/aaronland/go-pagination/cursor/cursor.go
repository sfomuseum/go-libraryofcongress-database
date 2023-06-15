// package cursor provides implementions of the pagintion.Options and pagination.Results interfaces for use with cursor or token-based pagination.
package cursor

import (
	"github.com/aaronland/go-pagination"
)

func NextCursor(r pagination.Results) string {

	if r.Method() != pagination.Cursor {
		return ""
	}

	return r.Next().(string)
}

func PreviousCursor(r pagination.Results) string {

	if r.Method() != pagination.Cursor {
		return ""
	}

	return r.Previous().(string)
}

func CursorFromOptions(opts pagination.Options) string {

	if opts.Method() != pagination.Cursor {
		return ""
	}

	return opts.Pointer().(string)
}
