# go-libraryofcongress-database

Go package providing simple database and server interfaces for the CSV files produced by the [sfomuseum/go-libraryofcongress](https://github.com/sfomuseum/go-libraryofcongress) package.

## Important

This is work in progress and not documented properly yet. The code will continue to move around in the short-term. Everything you see here is still in the "proof-of-concept" phase. It should work but may still have bugs and probably lacks features.

## Documentation

Documentation is incomplete at this time.

## Motivation

The first goal is to have a simple, bare-bones HTTP server for querying data in the CSV files produced by the [sfomuseum/go-libraryofcongress](https://github.com/sfomuseum/go-libraryofcongress) package.

The second goal is to be able to build, compile and deploy the web application and all its data (as SQLite databases) as a self-contained container image to a low-cost service like AWS App Runner or AWS Lambda Function URLs.

A third goal is to have a generic database interface such that the same code can be used with a variety of databases. As written the `server` tool only has a single database "driver" for querying SQLite databases but there are tools for indexing data in both Elasticsearch, OpenSeach, DynamoDB and SQLite databases.

## Databases

Databases implement the [LibraryOfCongressDatabase](LibraryOfCongressDatabase) interface which exposes two method: One for indexing source and another for querying them.

This package used to bundle a bunch of different implementations of the `LibraryOfCongressDatabase` interface but most of them have been moved in to separate packages. This package now only bundles two implementations:

* [SQLDatabase](sql/database.go) which implements the `LibraryOfCongressDatabase` interface using Go's `database/sql` abstraction layer. This package automatically loads the [modernc/sqlite](https://pkg.go.dev/modernc.org/sqlite) package for use with SQLite databases.

* [StdoutDatabase](stdout/database.go) which implements the indexing methods of the `LibraryOfCongressDatabase` interface and simply emits each record as a CSV-encoded row to STDOUT. This package is really only for testing and debugging purposes and to serve as a simple reference implementation for creating your own implementations.

Other database implementations include:

* https://github.com/sfomuseum/go-libraryofcongress-database-docstore
* https://github.com/sfomuseum/go-libraryofcongress-database-opensearch
* https://github.com/sfomuseum/go-libraryofcongress-database-elasticsearch
* https://github.com/sfomuseum/go-libraryofcongress-database-bleve

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" --tags fts5 -o bin/server cmd/server/main.go
go build -mod vendor -ldflags="-s -w" --tags fts5 -o bin/query cmd/query/main.go
go build -mod vendor -ldflags="-s -w" --tags fts5 -o bin/index cmd/index/main.go
```

### index

```
$> ./bin/index -h
  -database-uri string
    	A valid sfomuseum/go-libraryofcongress-database URI.
  -lcnaf-data string
    	The path to your LCNAF CSV data. If '-' then data will be read from STDIN.
  -lcsh-data string
    	The path to your LCSH CSV data. If '-' then data will be read from STDIN.
```

_Note that the index tool expects to index CSV data produced by the `sfomuseum/go-libraryofcongress` [parse-lcnaf](https://github.com/sfomuseum/go-libraryofcongress#parse-lcnaf) and [parse-lcsh](https://github.com/sfomuseum/go-libraryofcongress#parse-lcsh) tools._

#### SQL (SQLite)

```
$> ./bin/index -database-uri 'sql://sqlite?dsn=loc.db' -lcsh-data /usr/local/data/lcsh.csv.bz2
processed 5692 records in 1m0.001390571s (started 2021-10-27 15:52:35.790947 -0700 PDT m=+0.015128394)
processed 11161 records in 2m0.001847245s (started 2021-10-27 15:52:35.790947 -0700 PDT m=+0.015128394)
processed 16179 records in 3m0.000195064s (started 2021-10-27 15:52:35.790947 -0700 PDT m=+0.015128394)
processed 20693 records in 4m0.003592035s (started 2021-10-27 15:52:35.790947 -0700 PDT m=+0.015128394)
...time passes

processed 438053 records in 2h3m0.000624126s (started 2021-10-27 15:52:35.790947 -0700 PDT m=+0.015128394)
processed 441373 records in 2h4m0.002261248s (started 2021-10-27 15:52:35.790947 -0700 PDT m=+0.015128394)
processed 444805 records in 2h5m0.002327734s (started 2021-10-27 15:52:35.790947 -0700 PDT m=+0.015128394)
2021/10/27 17:58:09 Finished indexing lcsh

$> du -h -d 1 /usr/local/data/loc.db/
761M	loc.db/
```

#### STDIN

It is also possible to index data from `STDIN` by specifying the string "-" as the `-lcsh-data` or `-lcnaf-data` URI to read.

For example, this command will stream and parse the contents of `https://id.loc.gov/download/lcsh.both.ndjson.zip` (using the `parse-lcsh` tool in the [sfomuseum/go-libraryofcongress](https://github.com/sfomuseum/go-libraryofcongress#parse-lcsh) package) and index each subject header in a SQLite database called `'loc.db`

```
$> ./parse-lcsh https://id.loc.gov/download/lcsh.both.ndjson.zip | \
	./index \
	-database-uri 'sql://sqlite?dsn=loc.db' \
	-lcsh-data -
```

### server

The `server` tool is a simple web interface providing humans and robots, both, the ability to query a database.

```
$> ./bin/server -h
Usage of ./bin/server:
  -database-uri string
    	A valid sfomuseum/go-libraryofcongress-database URI.
  -per-page int
    	The number of results to return per page (default 20)
  -server-uri string
    	A valid aaronland/go-http-server URI. (default "http://localhost:8080")
```

To start the server you might do something like this:

```
$> ./bin/server -database-uri 'sql://sqlite?dsn=data/lcsh.db' -per-page 10
2021/10/18 13:11:24 Listening on http://localhost:8080
```

And then if you opened `http://localhost:8080/?q=River&page=2` in a web browser you'd see this:

![](docs/images/www.png)

There is also an API endpoint for querying the data as JSON:

```
$> curl -s 'http://localhost:8080/api/query?q=SQL' | jq
{
  "results": [
    {
      "id": "sh96008008",
      "label": "PL/SQL (Computer program language)",
      "source": "lcsh"
    },
    {
      "id": "sh86006628",
      "label": "SQL (Computer program language)",
      "source": "lcsh"
    },
    {
      "id": "sh90004874",
      "label": "SQL*PLUS (Computer program language)",
      "source": "lcsh"
    },
    {
      "id": "sh87001812",
      "label": "SQL/ORACLE (Computer program language)",
      "source": "lcsh"
    }
  ],
  "pagination": {
    "total": 4,
    "per_page": 10,
    "page": 1,
    "pages": 1,
    "next_page": 0,
    "previous_page": 0,
    "pages_range": []
  }
}
```

### query

```
$> ./bin/query -h
  -cursor-pagination
    	Signal that pagination is cursor-based rather than countable.
  -database-uri string
    	A valid sfomuseum/go-libraryofcongress-database URI.
```

The `query` tool is a command-line application to perform fulltext queries against a database generated using data produced by the tools in `sfomuseum/go-libraryofcongress` package.

#### SQL (SQLite)

```
$> ./bin/query \
	-database-uri 'sql://sqlite?dsn=test.db' \
	Montreal
	
lcsh:sh85087079 Montreal River (Ont.)
lcsh:sh2010014761 Alfa Romeo Montreal automobile
lcsh:sh2017003022 Montreal Massacre, Montréal, Québec, 1989
```

## Docker

Yes, there is a Dockerfile for the `server` tool. The simplest way to get started is to run the `docker` target in this package's Makefile:

```
$> make docker
```

And then to start the server:

```
$> docker run -it -p 8080:8080 \
	-e LIBRARYOFCONGRESS_DATABASE_URI='sql://sqlite?dsn=/usr/local/data/lcsh.db' \
	-e LIBRARYOFCONGRESS_SERVER_URI='http://0.0.0.0:8080' \
	libraryofcongress-server
```

And then visit `http://localhost:8080` in a web browser.

### Notes

* As written the `Dockerfile` will copy all files ending in `.db` in the `data` folder in to the container's `/usr/local/data` folder.

## See also

* https://github.com/sfomuseum/go-libraryofcongress
