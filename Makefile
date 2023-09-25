GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

cli:
	go build -mod $(GOMOD) -ldflags="-s -w" --tags fts5 -o bin/server cmd/server/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" --tags fts5 -o bin/query cmd/query/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" --tags fts5 -o bin/index cmd/index/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/create-dynamodb-csv cmd/create-dynamodb-csv/main.go

docker:
	docker build -t libraryofcongress-server .

# https://aws.amazon.com/about-aws/whats-new/2018/08/use-amazon-dynamodb-local-more-easily-with-the-new-docker-image/
# https://hub.docker.com/r/amazon/dynamodb-local/

dynamo-local:
	docker run --rm -it -p 8000:8000 amazon/dynamodb-local

dynamo-tables-local:
	go run -mod vendor cmd/create-dynamodb-tables/main.go \
		-refresh \
		-client-uri 'awsdynamodb://libraryofcongress?local=true&partition_key=Id'
