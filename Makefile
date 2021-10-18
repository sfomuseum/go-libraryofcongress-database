cli:
	go build -mod vendor -o bin/server cmd/server/main.go

docker:
	docker build -t libraryofcongress-server .
