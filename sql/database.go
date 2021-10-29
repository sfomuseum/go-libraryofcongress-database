package sql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aaronland/go-pagination"
	pg_sql "github.com/aaronland/go-pagination-sql"
	"github.com/sfomuseum/go-libraryofcongress-database"
	"github.com/sfomuseum/go-timings"
	"net/url"
)

type SQLDatabase struct {
	database.LibraryOfCongressDatabase
	db *sql.DB
}

func init() {
	ctx := context.Background()
	database.RegisterLibraryOfCongressDatabase(ctx, "sql", NewSQLDatabase)
}

func NewSQLDatabase(ctx context.Context, uri string) (database.LibraryOfCongressDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	// u.Scheme is assumed to be sql://
	engine := u.Host

	q := u.Query()
	dsn := q.Get("dsn")

	db, err := sql.Open(engine, dsn)

	if err != nil {
		return nil, fmt.Errorf("Failed to create database for '%s', %w", dsn, err)
	}

	err = db.Ping()

	if err != nil {
		return nil, fmt.Errorf("Failed to contact database with '%s', %w", dsn, err)
	}

	if engine == "sqlite3" {

		pragma := []string{
			"PRAGMA JOURNAL_MODE=OFF",
			"PRAGMA SYNCHRONOUS=OFF",
			"PRAGMA LOCKING_MODE=EXCLUSIVE",
			// https://www.gaia-gis.it/gaia-sins/spatialite-cookbook/html/system.html
			"PRAGMA PAGE_SIZE=4096",
			"PRAGMA CACHE_SIZE=1000000",
		}

		for _, p := range pragma {

			_, err = db.ExecContext(ctx, p)

			if err != nil {
				return nil, fmt.Errorf("Failed to set '%s', %w", p, err)
			}
		}
	}

	tables_sql := []string{
		"CREATE VIRTUAL TABLE IF NOT EXISTS search USING fts5(id, source, label)",
	}

	for _, q := range tables_sql {

		_, err = db.ExecContext(ctx, q)

		if err != nil {
			return nil, fmt.Errorf("Failed to execute tables SQL '%s', %w", q, err)
		}
	}

	sql_q := &SQLDatabase{
		db: db,
	}

	return sql_q, nil
}

func (sql_db *SQLDatabase) Index(ctx context.Context, sources []*database.Source, monitor timings.Monitor) error {

	for _, src := range sources {

		err := sql_db.indexSource(ctx, src, monitor)

		if err != nil {
			return fmt.Errorf("Failed to index %s, %w", src.Label, err)
		}
	}

	return nil
}

func (sql_db *SQLDatabase) indexSource(ctx context.Context, src *database.Source, monitor timings.Monitor) error {

	cb := func(ctx context.Context, row map[string]string) error {

		q := "INSERT OR REPLACE INTO search (id, source, label) VALUES (?, ?, ?	)"

		args := []interface{}{
			row["id"],
			row["source"],
			row["label"],
		}

		tx, err := sql_db.db.Begin()

		if err != nil {
			return fmt.Errorf("Failed to start transaction, %w", err)
		}

		stmt, err := tx.Prepare(q)

		if err != nil {
			return fmt.Errorf("Failed to prepare statement, %w", err)
		}

		defer stmt.Close()

		_, err = stmt.Exec(args...)

		if err != nil {
			return fmt.Errorf("Failed to execute query, %w", err)
		}

		err = tx.Commit()

		if err != nil {
			fmt.Errorf("Failed to commit transaction, %w", err)
		}

		go monitor.Signal(ctx)
		return nil
	}

	return src.Index(ctx, cb)
}

func (sql_db *SQLDatabase) Query(ctx context.Context, q string, pg_opts pagination.PaginationOptions) ([]*database.QueryResult, pagination.Pagination, error) {

	query_str := "SELECT id, label, source FROM search WHERE label MATCH ?  OR id MATCH ? ORDER BY label"

	pg_rsp, err := pg_sql.QueryPaginated(sql_db.db, pg_opts, query_str, q, q)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to query, %w", err)
	}

	results := make([]*database.QueryResult, 0)

	rows := pg_rsp.Rows()
	pagination := pg_rsp.Pagination()

	defer rows.Close()

	for rows.Next() {

		var id string
		var label string
		var source string

		err := rows.Scan(&id, &label, &source)

		if err != nil {
			return nil, nil, fmt.Errorf("Failed to scan row, %w", err)
		}

		r := &database.QueryResult{
			Id:     id,
			Label:  label,
			Source: source,
		}

		results = append(results, r)
	}

	err = rows.Close()

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to close database rows, %w", err)
	}

	err = rows.Err()

	if err != nil {
		return nil, nil, fmt.Errorf("Problem retrieving database rows, %w", err)
	}

	return results, pagination, nil
}
