module github.com/sfomuseum/go-libraryofcongress-database

go 1.18

// Pin to elastic/go-elasticsearch/v7 v7.13.0 because later versions
// don't work with AWS Elasticsearch anymore. Sigh...

require (
	github.com/aaronland/go-http-bootstrap v0.4.0
	github.com/aaronland/go-http-sanitize v0.0.8
	github.com/aaronland/go-http-server v1.2.0
	github.com/aaronland/go-pagination v0.2.0
	github.com/aaronland/go-pagination-sql v0.2.0
	github.com/aaronland/go-roster v1.0.0
	github.com/aaronland/go-sqlite v0.2.2
	github.com/blevesearch/bleve v1.0.14
	github.com/cenkalti/backoff/v4 v4.2.1
	github.com/elastic/go-elasticsearch/v7 v7.13.0
	github.com/mattn/go-sqlite3 v1.14.17
	github.com/sfomuseum/go-csvdict v1.0.0
	github.com/sfomuseum/go-flags v0.10.0
	github.com/sfomuseum/go-timings v1.2.1
)

require (
	github.com/RoaringBitmap/roaring v0.4.23 // indirect
	github.com/aaronland/go-http-rewrite v1.1.0 // indirect
	github.com/aaronland/go-http-static v0.0.3 // indirect
	github.com/aaronland/go-log/v2 v2.0.0 // indirect
	github.com/akrylysov/algnhsa v1.0.0 // indirect
	github.com/aws/aws-lambda-go v1.37.0 // indirect
	github.com/blevesearch/go-porterstemmer v1.0.3 // indirect
	github.com/blevesearch/mmap-go v1.0.2 // indirect
	github.com/blevesearch/segment v0.9.0 // indirect
	github.com/blevesearch/snowballstem v0.9.0 // indirect
	github.com/blevesearch/zap/v11 v11.0.14 // indirect
	github.com/blevesearch/zap/v12 v12.0.14 // indirect
	github.com/blevesearch/zap/v13 v13.0.6 // indirect
	github.com/blevesearch/zap/v14 v14.0.5 // indirect
	github.com/blevesearch/zap/v15 v15.0.3 // indirect
	github.com/couchbase/vellum v1.0.2 // indirect
	github.com/glycerine/go-unsnap-stream v0.0.0-20181221182339-f9677308dec2 // indirect
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/jtacoma/uritemplates v1.0.0 // indirect
	github.com/mschoch/smat v0.2.0 // indirect
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/sfomuseum/go-http-rollup v0.0.2 // indirect
	github.com/sfomuseum/iso8601duration v1.1.0 // indirect
	github.com/steveyen/gtreap v0.1.0 // indirect
	github.com/tdewolff/minify/v2 v2.12.4 // indirect
	github.com/tdewolff/parse/v2 v2.6.4 // indirect
	github.com/tinylib/msgp v1.1.0 // indirect
	github.com/whosonfirst/go-sanitize v0.1.0 // indirect
	github.com/willf/bitset v1.1.10 // indirect
	go.etcd.io/bbolt v1.3.5 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
)
