package tables

import (
	"context"
	_ "embed"
	"fmt"
	
	"github.com/aaronland/go-sqlite"
)

//go:embed identifiers.schema
var identifiers_schema string

// type IdentifiersTable implements the `sqlite.Table` interface for mapping LoC identifiers to their corresponding label.
type IdentifiersTable struct {
	sqlite.Table
	// The name of the SQLite database table containing identifiers.
	name string
}

//  NewIdentifiersTableWithTable() returns a new `IdentifiersTable` instance for use with the database identifier by 'db'.
func NewIdentifiersTableWithDatabase(ctx context.Context, db sqlite.Database) (sqlite.Table, error) {

	t, err := NewIdentifiersTable(ctx)

	if err != nil {
		return nil, fmt.Errorf("Failed to create identifiers table, %w", err)
	}

	err = t.InitializeTable(ctx, db)

	if err != nil {
		return nil, fmt.Errorf("Failed to initialize identifiers table, %w", err)
	}

	return t, nil
}

// NewIdentifiersTableWithTable() returns a new `IdentifiersTable` instance.
func NewIdentifiersTable(ctx context.Context) (sqlite.Table, error) {

	t := IdentifiersTable{
		name: "identifiers",
	}

	return &t, nil
}

// InitializeTable() will ensure the an identifiers table has been created in the database represented by 'db'.
func (t *IdentifiersTable) InitializeTable(ctx context.Context, db sqlite.Database) error {
	return sqlite.CreateTableIfNecessary(ctx, db, t)
}

// Name() returns the name of the identifiers table.
func (t *IdentifiersTable) Name() string {
	return t.name
}

// Schema() returns the schema used to create the identifiers table.
func (t *IdentifiersTable) Schema() string {
	return identifiers_schema
}

// IndexRecord() indexes 'i' in the database represented by 'db'.
func (t *IdentifiersTable) IndexRecord(ctx context.Context, db sqlite.Database, i interface{}) error {
	return t.IndexRow(ctx, db, i.(map[string]string))
}

// IndexRecord() indexes 'row' in the database represented by 'db'.
func (t *IdentifiersTable) IndexRow(ctx context.Context, db sqlite.Database, row map[string]string) error {

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
		return fmt.Errorf("Failed to connect to database, %w", err)
	}

	tx, err := conn.Begin()

	if err != nil {
		return fmt.Errorf("Failed to begin transaction, %w", err)
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
