# go-libraryofcongress-database

Go package providing simple database and server interfaces for the CSV files produced by the [sfomuseum/go-libraryofcongress](https://github.com/sfomuseum/go-libraryofcongress) package.

## Important

This is work in progress and not documented properly yet. The code will continue to move around in the short-term. Everything you see here is still in the "proof-of-concept" phase. It should work but may still have bugs and probably lacks features.

## Motivation

The first goal is to have a simple, bare-bones HTTP server for querying data in the CSV files produced by the [sfomuseum/go-libraryofcongress](https://github.com/sfomuseum/go-libraryofcongress) package.

The second goal is to be able to build, compile and deploy the web application and all its data (as SQLite databases) as a self-contained container image to a low-cost service like AWS App Runner.

A third goal is to have a generic database interface such that the same code can be used with a variety of databases. As written the `server` tool only has a single database "driver" for querying SQLite databases but there are tools for indexing data in both Elasticsearch and SQLite databases.

## Data

A sample SQLite database for Library of Congress subject headings is currently included with this package in the [data folder](data). Some notes:

* This database is stored using `git-lfs`.
* This databases was created using the `to-sqlite` tool described below.
* It is not clear whether an equivalent (or combined) database for Library of Congress named authorities will ever be included because it is very large.
* Eventually bundled data may be removed entirely.
* As written the code only handles a subset of all the possible (CSV) columns produced by the `sfomuseum/go-libraryofcongress` tools. Specifically: `id` and `label`. A third `source` column is appended to the databases to distinguish between Library of Congress subject heading and name authority file records.

## Databases

_To be written_

## Tools

### server

The `server` tool is a simple web interface providing humans and robots, both, the ability to query a database.

```
$> ./bin/server -h
Usage of ./bin/server:
  -database-uri string
    	A valid sfomuseum/go-libraryofcongress-database URI. (default "sql://sqlite3?dsn=data/lcsh.db")
  -per-page int
    	The number of results to return per page (default 20)
  -server-uri string
    	A valid aaronland/go-http-server URI. (default "http://localhost:8080")
```

To start the server you might do something like this:

```
$> ./bin/server -database-uri 'sql://sqlite3?dsn=data/lcsh.db' -per-page 10
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

#### Notes

* The `server` tool only supports SQLite databases and Bleve indexes as of this writing. See below for details on tools to produce these data sources.
* The `server` tool does not yet have the ability to define custom prefixes for URLs. For the time being it is assumed that everything is served from a root `/` URL.

### query

```
$> ./bin/query -database-uri bleve:///usr/local/data/libraryofcongress.db Montreal
lcsh:sh85087079 Montreal River (Ont.)
lcsh:sh2010014761 Alfa Romeo Montreal automobile
lcsh:sh2017003022 Montreal Massacre, Montréal, Québec, 1989
```

### to-bleve

The `to-bleve` tool will index CSV data produced by the tools in `sfomuseum/go-libraryofcongress` in an Bleve index.

```
> ./bin/to-bleve -h
Usage of ./bin/to-bleve:
  -index string
    	The path to the Bleve index you want to create. (default "libraryofcongress.db")
  -lcnaf-data string
    	The path to your LCNAF CSV data.
  -lcsh-data string
    	The path to your LCSH CSV data.
```

```
$> ./bin/to-bleve -index /usr/local/data/libraryofcongress.db -lcsh-data /usr/local/data/lcsh.csv.bz2
processed 5692 records in 1m0.001390571s (started 2021-10-27 15:52:35.790947 -0700 PDT m=+0.015128394)
processed 11161 records in 2m0.001847245s (started 2021-10-27 15:52:35.790947 -0700 PDT m=+0.015128394)
processed 16179 records in 3m0.000195064s (started 2021-10-27 15:52:35.790947 -0700 PDT m=+0.015128394)
processed 20693 records in 4m0.003592035s (started 2021-10-27 15:52:35.790947 -0700 PDT m=+0.015128394)
...time passes

processed 438053 records in 2h3m0.000624126s (started 2021-10-27 15:52:35.790947 -0700 PDT m=+0.015128394)
processed 441373 records in 2h4m0.002261248s (started 2021-10-27 15:52:35.790947 -0700 PDT m=+0.015128394)
processed 444805 records in 2h5m0.002327734s (started 2021-10-27 15:52:35.790947 -0700 PDT m=+0.015128394)
2021/10/27 17:58:09 Finished indexing lcsh

$> du -h -d 1 /usr/local/data/libraryofcongress.db/
761M	libraryofcongress.db/
```

And then you can use the `query` tool (described above) to query the database:

```
$> ./bin/query -database-uri bleve:///usr/local/data/libraryofcongress.db Montreal
lcsh:sh85087079 Montreal River (Ont.)
lcsh:sh2010014761 Alfa Romeo Montreal automobile
lcsh:sh2017003022 Montreal Massacre, Montréal, Québec, 1989
```

### to-elasticsearch

The `to-elasticsearch` tool will index CSV data produced by the tools in `sfomuseum/go-libraryofcongress` in an Elasticsearch index.

```
$> ./bin/to-elasticsearch -h /Users/asc/sfomuseum/go-libraryofcongress-database                                                    
Usage of ./bin/to-elasticsearch:
  -elasticsearch-endpoint string
    	The Elasticsearch endpoint where data should be indexed. (default "http://localhost:9200")
  -elasticsearch-index string
    	The Elasticsearch index where data should be stored. (default "libraryofcongress")
  -lcnaf-data string
    	The path to your LCNAF CSV data.
  -lcsh-data string
    	The path to your LCSH CSV data.
  -workers int
    	The number of concurrent workers to use when indexing data. (default 10)
```

### to-sqlite

The `to-sqlite` tool will index CSV data produced by the tools in `sfomuseum/go-libraryofcongress` in a SQLite database.

```
$> ./bin/to-sqlite -h
Usage of ./bin/to-sqlite:
  -dsn string
    	The SQLite DSN for the database you want to create. (default "libraryofcongress.db")
  -lcnaf-data string
    	The path to your LCNAF CSV data.
  -lcsh-data string
    	The path to your LCSH CSV data.
```

## Important

This works. But it is very slow to index data. For example, here are some timings from an attempt to index both the Library of Congress subject headings and name authority file in a single SQLite database:

```
processed 253390 records in 4h13m0.03018635s (started 2021-10-19 17:58:06.533669374 +0000 UTC m=+0.002113478)
```

So it takes between 16-20  hours to index the Library of Congress subject headings. Here's what happens when the same database is used to index the Library of Congress Name Authority File as well:

```
processed 614560 records in 22h12m0.000038787s (started 2021-10-19 17:58:06.533669374 +0000 UTC m=+0.002113478)
processed 713753 records in 29h52m0.028314524s (started 2021-10-19 17:58:06.533669374 +0000 UTC m=+0.002113478)
processed 889205 records in 46h20m0.014818217s (started 2021-10-19 17:58:06.533669374 +0000 UTC m=+0.002113478)
processed 954607 records in 53h26m0.040308808s (started 2021-10-19 17:58:06.533669374 +0000 UTC m=+0.002113478)
processed 1152145 records in 78h4m0.012184985s (started 2021-10-19 17:58:06.533669374 +0000 UTC m=+0.002113478)
processed 1489494 records in 131h5m0.019974672s (started 2021-10-19 17:58:06.533669374 +0000 UTC m=+0.002113478)
processed 1556025 records in 143h10m0.040390679s (started 2021-10-19 17:58:06.533669374 +0000 UTC m=+0.002113478)
processed 1704642 records in 172h9m0.014664052s (started 2021-10-19 17:58:06.533669374 +0000 UTC m=+0.002113478)
```

Remember, there are 11 million records in the Name Authority File so it's pretty easy to extropolate that the amount of time it will take to index another 10 million records will stretch in to weeks and probably months.

Indexing in SQLite is slow enough that any of the other alternatives may be preferable. For example it takes Elasticsearch a couple of minutes to index all 11 million records in a couple of minutes.

## Docker

Yes, there is a Dockerfile for the `server` tool. The simplest way to get started is to run the `docker` target in this package's Makefile:

```
$> make docker
```

And then to start the server:

```
$> docker run -it -p 8080:8080 \
	-e LIBRARYOFCONGRESS_DATABASE_URI='sql://sqlite3?dsn=/usr/local/data/lcsh.db' \
	-e LIBRARYOFCONGRESS_SERVER_URI='http://0.0.0.0:8080' \
	libraryofcongress-server
```

And then visit `http://localhost:8080` in a web browser.

### Notes

* As written the `Dockerfile` will copy all files ending in `.db` in the `data` folder in to the container's `/usr/local/data` folder.

## See also

* https://github.com/sfomuseum/go-libraryofcongress