cli:
	go build -mod vendor -o bin/server cmd/server/main.go
	go build -mod vendor -o bin/to-bleve cmd/to-bleve/main.go
	go build -mod vendor -o bin/to-elasticsearch cmd/to-elasticsearch/main.go
	go build -mod vendor --tags fts5 -o bin/to-sqlite cmd/to-sqlite/main.go

docker:
	docker build -t libraryofcongress-server .
