package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aaronland/go-pagination"
	pg_sql "github.com/aaronland/go-pagination-sql"
	"net/url"
)

type SQLDatabase struct {
	LibraryOfCongressDatabase
	db *sql.DB
}

func init() {
	ctx := context.Background()
	RegisterLibraryOfCongressDatabase(ctx, "sql", NewSQLDatabase)
}

func NewSQLDatabase(ctx context.Context, uri string) (LibraryOfCongressDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	// u.Scheme is assumed to be sql://
	driver := u.Host

	q := u.Query()
	dsn := q.Get("dsn")

	db, err := sql.Open(driver, dsn)

	if err != nil {
		return nil, fmt.Errorf("Failed to create database for '%s', %w", dsn, err)
	}

	err = db.Ping()

	if err != nil {
		return nil, fmt.Errorf("Failed to contact database with '%s', %w", dsn, err)
	}

	sql_q := &SQLDatabase{
		db: db,
	}

	return sql_q, nil
}

func (sql_q *SQLDatabase) Query(ctx context.Context, q string, pg_opts pagination.PaginationOptions) ([]*QueryResult, pagination.Pagination, error) {

	query_str := "SELECT id, label, source FROM search WHERE label MATCH ?  OR id MATCH ? ORDER BY label"

	pg_rsp, err := pg_sql.QueryPaginated(sql_q.db, pg_opts, query_str, q, q)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to query, %w", err)
	}

	results := make([]*QueryResult, 0)

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

		r := &QueryResult{
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
