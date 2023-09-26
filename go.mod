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
	github.com/sfomuseum/go-csvdict v1.0.0
	github.com/sfomuseum/go-flags v0.10.0
	github.com/sfomuseum/go-timings v1.2.1
)

require (
	github.com/aaronland/go-http-rewrite v1.1.0 // indirect
	github.com/aaronland/go-http-static v0.0.3 // indirect
	github.com/aaronland/go-log/v2 v2.0.0 // indirect
	github.com/akrylysov/algnhsa v1.0.0 // indirect
	github.com/aws/aws-lambda-go v1.37.0 // indirect
	github.com/jtacoma/uritemplates v1.0.0 // indirect
	github.com/sfomuseum/go-http-rollup v0.0.2 // indirect
	github.com/sfomuseum/iso8601duration v1.1.0 // indirect
	github.com/tdewolff/minify/v2 v2.12.4 // indirect
	github.com/tdewolff/parse/v2 v2.6.4 // indirect
	github.com/whosonfirst/go-sanitize v0.1.0 // indirect
	golang.org/x/net v0.13.0 // indirect
)
