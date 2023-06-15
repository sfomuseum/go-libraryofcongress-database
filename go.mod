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
	github.com/aaronland/gocloud-docstore v0.0.5
	github.com/blevesearch/bleve v1.0.14
	github.com/cenkalti/backoff/v4 v4.2.1
	github.com/elastic/go-elasticsearch/v7 v7.13.0
	github.com/mattn/go-sqlite3 v1.14.17
	github.com/sfomuseum/go-csvdict v1.0.0
	github.com/sfomuseum/go-flags v0.10.0
	github.com/sfomuseum/go-timings v1.2.1
	gocloud.dev v0.29.0
)

require (
	github.com/RoaringBitmap/roaring v0.4.23 // indirect
	github.com/aaronland/go-aws-dynamodb v0.0.5 // indirect
	github.com/aaronland/go-aws-session v0.0.6 // indirect
	github.com/aaronland/go-http-rewrite v1.1.0 // indirect
	github.com/aaronland/go-http-static v0.0.3 // indirect
	github.com/aaronland/go-log/v2 v2.0.0 // indirect
	github.com/aaronland/go-string v0.1.2 // indirect
	github.com/akrylysov/algnhsa v1.0.0 // indirect
	github.com/aws/aws-lambda-go v1.37.0 // indirect
	github.com/aws/aws-sdk-go v1.44.200 // indirect
	github.com/aws/aws-sdk-go-v2 v1.17.4 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.18.12 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.13.12 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.28 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.29 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.22 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.12.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.14.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.18.3 // indirect
	github.com/aws/smithy-go v1.13.5 // indirect
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
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/google/wire v0.5.0 // indirect
	github.com/googleapis/gax-go/v2 v2.7.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
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
	github.com/willf/bitset v1.1.11 // indirect
	go.etcd.io/bbolt v1.3.6 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.110.0 // indirect
	google.golang.org/genproto v0.0.0-20230209215440-0dfe4f8abfcc // indirect
	google.golang.org/grpc v1.53.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
