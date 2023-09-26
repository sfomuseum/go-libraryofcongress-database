GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

cli:
	go build -mod $(GOMOD) -ldflags="-s -w" --tags fts5 -o bin/server cmd/server/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" --tags fts5 -o bin/query cmd/query/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" --tags fts5 -o bin/index cmd/index/main.go
