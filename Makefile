cli:
	go build -mod vendor --tags fts5 -o bin/server cmd/server/main.go
	go build -mod vendor --tags fts5 -o bin/query cmd/query/main.go
	go build -mod vendor --tags fts5 -o bin/index cmd/index/main.go

docker:
	docker build -t libraryofcongress-server .
