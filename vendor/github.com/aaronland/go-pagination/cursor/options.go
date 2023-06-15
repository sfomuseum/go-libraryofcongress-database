package cursor

import (
	"github.com/aaronland/go-pagination"
)

const PER_PAGE int64 = 10
const PAGE int64 = 1
const SPILL int64 = 2

type CursorOptions struct {
	pagination.Options
	perpage int64
	cursor  string
	column  string
}

func NewCursorOptions() (pagination.Options, error) {

	opts := &CursorOptions{
		perpage: PER_PAGE,
	}

	return opts, nil
}

func (opts *CursorOptions) Method() pagination.Method {
	return pagination.Cursor
}

func (opts *CursorOptions) Pointer(args ...interface{}) interface{} {

	if len(args) >= 1 {
		opts.cursor = args[0].(string)
	}

	return opts.cursor
}

func (opts *CursorOptions) PerPage(args ...int64) int64 {

	if len(args) >= 1 {
		opts.perpage = args[0]
	}

	return opts.perpage
}

func (opts *CursorOptions) Page(args ...int64) int64 {
	return 0
}

func (opts *CursorOptions) Spill(args ...int64) int64 {
	return 0
}

func (opts *CursorOptions) Column(args ...string) string {

	if len(args) >= 1 {
		opts.column = args[0]
	}

	return opts.column
}
