package countable

import (
	"github.com/aaronland/go-pagination"
)

const PER_PAGE int64 = 10
const PAGE int64 = 1
const SPILL int64 = 2
const COUNTABLE string = "*"

type CountablePaginationOptions struct {
	pagination.PaginationOptions
	perpage int64
	page    int64
	spill   int64
	column  string
}

func NewCountablePaginationOptions() (pagination.PaginationOptions, error) {

	opts := &CountablePaginationOptions{
		perpage: PER_PAGE,
		page:    PAGE,
		spill:   SPILL,
		column:  COUNTABLE,
	}

	return opts, nil
}

func (opts *CountablePaginationOptions) PerPage(args ...int64) int64 {

	if len(args) >= 1 {
		opts.perpage = args[0]
	}

	return opts.perpage
}

func (opts *CountablePaginationOptions) Page(args ...int64) int64 {

	if len(args) >= 1 {
		opts.page = args[0]
	}

	return opts.page
}

func (opts *CountablePaginationOptions) Spill(args ...int64) int64 {

	if len(args) >= 1 {
		opts.spill = args[0]
	}

	return opts.spill
}

func (opts *CountablePaginationOptions) Column(args ...string) string {

	if len(args) >= 1 {
		opts.column = args[0]
	}

	return opts.column
}
