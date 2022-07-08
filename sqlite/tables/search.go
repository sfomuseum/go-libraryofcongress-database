package tables

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/aaronland/go-sqlite"
)

//go:embed search.schema
var search_schema string

// type SeachTable implements the `sqlite.Table` interface for searching LoC data (using the SQLite FTS5 plugin).
type SearchTable struct {
	sqlite.Table
	// The name of the SQLite database table containing searchable records.
	name string
}

//  NewSearchTableWithTable() returns a new `SearchTable` instance for use with the database identifier by 'db'.
func NewSearchTableWithDatabase(ctx context.Context, db sqlite.Database) (sqlite.Table, error) {

	t, err := NewSearchTable(ctx)

	if err != nil {
		return nil, fmt.Errorf("Failed to create search table, %w", err)
	}

	err = t.InitializeTable(ctx, db)

	if err != nil {
		return nil, fmt.Errorf("Failed to initialize search table, %w", err)
	}

	return t, nil
}

// NewSearchTableWithTable() returns a new `SearchTable` instance.
func NewSearchTable(ctx context.Context) (sqlite.Table, error) {

	t := SearchTable{
		name: "search",
	}

	return &t, nil
}

// InitializeTable() will ensure the an search table has been created in the database represented by 'db'.
func (t *SearchTable) InitializeTable(ctx context.Context, db sqlite.Database) error {

	return sqlite.CreateTableIfNecessary(ctx, db, t)
}

// Name() returns the name of the search table.
func (t *SearchTable) Name() string {
	return t.name
}

// Schema() returns the schema used to create the identifiers table.
func (t *SearchTable) Schema() string {
	return search_schema
}

// IndexRecord() indexes 'i' in the database represented by 'db'.
func (t *SearchTable) IndexRecord(ctx context.Context, db sqlite.Database, i interface{}) error {
	return t.IndexRow(ctx, db, i.(map[string]string))
}

// IndexRecord() indexes 'row' in the database represented by 'db'.
func (t *SearchTable) IndexRow(ctx context.Context, db sqlite.Database, row map[string]string) error {

	sql := fmt.Sprintf(`INSERT OR REPLACE INTO %s (
		id, source, label
		) VALUES (
		?, ?, ?
		)`, t.Name()) // ON CONFLICT DO BLAH BLAH BLAH

	args := []interface{}{
		row["id"],
		row["source"],
		row["label"],
	}

	conn, err := db.Conn()

	if err != nil {
		return fmt.Errorf("Failed to create database connection, %w", err)
	}

	tx, err := conn.Begin()

	if err != nil {
		return fmt.Errorf("Failed to being transaction, %w", err)
	}

	s, err := tx.Prepare(fmt.Sprintf("DELETE FROM %s WHERE id = ?", t.Name()))

	if err != nil {
		return fmt.Errorf("Failed to prepare statement, %w", err)
	}

	defer s.Close()

	_, err = s.Exec(row["id"])

	if err != nil {
		return fmt.Errorf("Failed to delete rows for ID %s, %w", row["id"], err)
	}

	stmt, err := tx.Prepare(sql)

	if err != nil {
		return fmt.Errorf("Failed to prepare statement, %w", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(args...)

	if err != nil {
		return fmt.Errorf("Failed to execute statement, %w", err)
	}

	err = tx.Commit()

	if err != nil {
		return fmt.Errorf("Failed to commit transaction, %w", err)
	}

	return nil
}
