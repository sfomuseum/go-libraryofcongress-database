package pagination

import (
	"math"
	"net/url"
)

type PaginationOptions interface {
	PerPage(...int64) int64
	Page(...int64) int64
	Spill(...int64) int64
	Column(...string) string
}

type Pagination interface {
	Total() int64
	PerPage() int64
	Page() int64
	Pages() int64
	NextPage() int64
	PreviousPage() int64
	NextURL(u *url.URL) string
	PreviousURL(u *url.URL) string
	Range() []int64
}

func PagesForCount(opts PaginationOptions, total_count int64) int64 {

	per_page := int64(math.Max(1.0, float64(opts.PerPage())))
	spill := int64(math.Max(1.0, float64(opts.Spill())))

	if spill >= per_page {
		spill = per_page - 1
	}

	pages := int64(math.Ceil(float64(total_count) / float64(per_page)))
	return pages
}
