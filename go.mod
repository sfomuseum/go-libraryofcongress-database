module github.com/sfomuseum/go-libraryofcongress-database

go 1.16

// Pin to elastic/go-elasticsearch/v7 v7.13.0 because later versions
// don't work with AWS Elasticsearch anymore. Sigh...

require (
	github.com/aaronland/go-http-bootstrap v0.1.0
	github.com/aaronland/go-http-sanitize v0.0.6
	github.com/aaronland/go-http-server v1.0.0
	github.com/aaronland/go-pagination v0.2.0
	github.com/aaronland/go-pagination-sql v0.2.0
	github.com/aaronland/go-roster v1.0.0
	github.com/aaronland/go-sqlite v0.2.0
	github.com/blevesearch/bleve v1.0.14
	github.com/cenkalti/backoff/v4 v4.1.3
	github.com/elastic/go-elasticsearch/v7 v7.13.0
	github.com/mattn/go-sqlite3 v1.14.13
	github.com/sfomuseum/go-csvdict v1.0.0
	github.com/sfomuseum/go-flags v0.8.2
	github.com/sfomuseum/go-timings v1.0.0
	golang.org/x/sys v0.0.0-20220708085239-5a0f0661e09d // indirect
)
