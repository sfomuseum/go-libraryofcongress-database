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

_TBW_

### Bleve

* https://blevesearch.com/

### Elasticsearch

* https://www.elastic.co/elastic-stack/

### SQlite

* https://sqlite.org/

## Database URIs

### bleve

```
bleve://{PATH_TO_DATABASE}
```

For example:

```
bleve:///usr/local/data/loc.db
```

### elasticsearch

```
elasticsearch://?endpoint={ELASTICSEARCH_INDEX}&index={ELASTICSEARCH_INDEX}
```

For example:

```
elasticsearch://?endpoint=http://localhost:9200&index=loc
```

### sqlite

```
sql://{ENGINE}?dsn={SQLITE_DSN}
```

For example:

```
sql://sqlite3?dsn=/usr/local/data/loc.db
```

## Tools

### index

#### bleve

```
$> ./bin/index -database-uri bleve:///usr/local/data/libraryofcongress.db -lcsh-data /usr/local/data/lcsh.csv.bz2
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

#### elasticsearch

_TBW_

#### sqlite

```
$> time ./bin/index -database-uri 'sql://sqlite3?dsn=test.db' -lcsh-data ~/Desktop/lcsh.csv.bz2
33.169u 15.511s 0:36.01 135.1%	0+0k 0+0io 0pf+0w

$> du -h test.db 
 46M	test.db
```

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

The `query` tool is a command-line application to perform fulltext queries against a database generated using data produced by the tools in `sfomuseum/go-libraryofcongress` package.

#### bleve

```
$> ./bin/query -database-uri bleve:///usr/local/data/libraryofcongress.db Montreal
lcsh:sh85087079 Montreal River (Ont.)
lcsh:sh2010014761 Alfa Romeo Montreal automobile
lcsh:sh2017003022 Montreal Massacre, Montr??al, Qu??bec, 1989
```

#### sqlite

```
$> ./bin/query -database-uri 'sql://sqlite3?dsn=test.db' Montreal
lcsh:sh2010014761 Alfa Romeo Montreal automobile
lcsh:sh94006536 Boulevard Saint-Laurent (Montr??al, Qu??bec)
lcsh:sh2009118684 Central business districts--Qu??bec (Province)--Montr??al--Maps
lcsh:sh2008002760 Fur Trade at Lachine National Historic Site (Montr??al, Qu??bec)
lcsh:sh85073824 Lachine Canal (Montr??al, Qu??bec)
lcsh:sh2008003685 Lachine Canal National Historic Site (Montr??al, Qu??bec)
lcsh:sh2008002035 Louis-Joseph Papineau National Historic Site (Montr??al, Qu??bec)
lcsh:sh86005383 Maison Saint-Gabriel (Montr??al, Qu??bec)
lcsh:sh2008107459 Marriage records--Qu??bec (Province)--Montr??al
lcsh:sh2017003022 Montreal Massacre, Montr??al, Qu??bec, 1989
lcsh:sh85087079 Montreal River (Ont.)
lcsh:sh2008115936 Montr??al (Qu??bec)--Biography
lcsh:sh2008115937 Montr??al (Qu??bec)--Fiction
lcsh:sh2008107460 Montr??al (Qu??bec)--Genealogy
lcsh:sh2008115938 Montr??al (Qu??bec)--Guidebooks
lcsh:sh95002319 Montr??al (Qu??bec)--History
lcsh:sh95002320 Montr??al (Qu??bec)--History--Siege, 1775
lcsh:sh85087078 Montr??al Island (Qu??bec : Island)
lcsh:sh2007000305 Parc Belmont (Montr??al, Qu??bec)
lcsh:sh2002003622 Parc Jarry (Montr??al, Qu??bec)
lcsh:sh88006626 Parc Sohmer (Montr??al, Qu??bec)
lcsh:sh2002003824 Parc du Mont-Royal (Montr??al, Qu??bec)
lcsh:sh2005005587 Place Ville Marie (Montr??al, Qu??bec)
lcsh:sh93002004 Place d'Armes (Montr??al, Qu??bec)
lcsh:sh93002901 Pointe-??-Calli??re Site (Montr??al, Qu??bec)
lcsh:sh94001599 Pont Victoria (Montr??al, Qu??bec)
lcsh:sh2001005269 Poudri??re (Montr??al, Qu??bec)
lcsh:sh2015001674 Rue Sainte-Catherine (Montr??al, Qu??bec)
lcsh:sh2005005292 Rue Sainte-H??l??ne (Montr??al, Qu??bec)
lcsh:sh2007003084 Sherbrooke Street (Montr??al, Qu??bec)
lcsh:sh2008002166 Sir George-??tienne Cartier National Historic Site (Montr??al, Qu??bec)
lcsh:sh2002005306 Stock Exchange Tower (Montr??al, Qu??bec)
```

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