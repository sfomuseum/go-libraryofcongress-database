package www

import (
	"fmt"
	"github.com/aaronland/go-http-sanitize"
	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/countable"
	"github.com/sfomuseum/go-libraryofcongress-database"
	_ "log"
	"net/http"
)

type QueryVars struct {
	Query      string
	Results    []*database.QueryResult
	Pagination pagination.Pagination
	Error      error
}

func QueryHandler(opts *Options) (http.Handler, error) {

	t := opts.Templates.Lookup("query")

	if t == nil {
		return nil, fmt.Errorf("Missing query template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		vars := QueryVars{}

		q, err := sanitize.GetString(req, "q")

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		if q == "" {
			renderTemplate(rsp, t, vars)
			return
		}

		vars.Query = q

		pg_opts, err := countable.NewCountablePaginationOptions()

		if err != nil {
			vars.Error = err
			renderTemplate(rsp, t, vars)
			return
		}

		pg_opts.PerPage(opts.PerPage)

		page, err := sanitize.GetInt64(req, "page")

		if err != nil {
			vars.Error = err
			renderTemplate(rsp, t, vars)
			return
		}

		if page > 0 {
			pg_opts.Page(page)
		}

		results, pagination, err := opts.Database.Query(ctx, q, pg_opts)

		if err != nil {
			vars.Error = err
			renderTemplate(rsp, t, vars)
			return
		}

		vars.Results = results
		vars.Pagination = pagination

		renderTemplate(rsp, t, vars)
		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
