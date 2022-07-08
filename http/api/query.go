package api

import (
	"encoding/json"
	"github.com/aaronland/go-http-sanitize"
	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/countable"
	"github.com/sfomuseum/go-libraryofcongress-database"
	_ "log"
	"net/http"
)

type QueryResponse struct {
	Results    []*database.QueryResult `json:"results"`
	Pagination pagination.Results      `json:"pagination"`
}

func QueryHandler(opts *Options) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		pg_opts, err := countable.NewCountableOptions()

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		pg_opts.PerPage(opts.PerPage)

		q, err := sanitize.GetString(req, "q")

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		if q == "" {
			http.Error(rsp, "Missing query (q) parameter", http.StatusBadRequest)
			return
		}

		page, err := sanitize.GetInt64(req, "page")

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		if page > 0 {
			pg_opts.Pointer(page)
		}

		results, pagination, err := opts.Database.Query(ctx, q, pg_opts)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		query_rsp := QueryResponse{
			Results:    results,
			Pagination: pagination,
		}

		rsp.Header().Set("Content-type", "application/json")

		enc := json.NewEncoder(rsp)
		err = enc.Encode(query_rsp)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
