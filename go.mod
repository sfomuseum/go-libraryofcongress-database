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
	modernc.org/sqlite v1.25.0
)

require (
	github.com/aaronland/go-http-rewrite v1.1.0 // indirect
	github.com/aaronland/go-http-static v0.0.3 // indirect
	github.com/aaronland/go-log/v2 v2.0.0 // indirect
	github.com/akrylysov/algnhsa v1.0.0 // indirect
	github.com/aws/aws-lambda-go v1.37.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jtacoma/uritemplates v1.0.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/sfomuseum/go-http-rollup v0.0.2 // indirect
	github.com/sfomuseum/iso8601duration v1.1.0 // indirect
	github.com/tdewolff/minify/v2 v2.12.4 // indirect
	github.com/tdewolff/parse/v2 v2.6.4 // indirect
	github.com/whosonfirst/go-sanitize v0.1.0 // indirect
	golang.org/x/mod v0.3.0 // indirect
	golang.org/x/net v0.13.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/tools v0.0.0-20201124115921-2c860bdd6e78 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	lukechampine.com/uint128 v1.2.0 // indirect
	modernc.org/cc/v3 v3.40.0 // indirect
	modernc.org/ccgo/v3 v3.16.13 // indirect
	modernc.org/libc v1.24.1 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.6.0 // indirect
	modernc.org/opt v0.1.3 // indirect
	modernc.org/strutil v1.1.3 // indirect
	modernc.org/token v1.0.1 // indirect
)
