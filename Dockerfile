FROM golang:1.17-alpine as builder

RUN mkdir /build

COPY . /build/go-libraryofcongress-database

RUN apk update && apk upgrade \
    && apk add make libc-dev gcc git \
    && cd /build/go-libraryofcongress-database \
    && go build -mod vendor -o /usr/local/bin/libraryofcongress-server cmd/server/main.go    

FROM alpine:latest

COPY --from=builder /usr/local/bin/libraryofcongress-server /usr/local/bin/

RUN mkdir /usr/local/data

COPY data/*.db /usr/local/data/

RUN apk update && apk upgrade \
    && apk add ca-certificates

ENTRYPOINT ["/usr/local/bin/libraryofcongress-server"]