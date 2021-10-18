package tables

import (
	"context"
	"fmt"
	"github.com/aaronland/go-sqlite"
	_ "log"
)

type SearchTable struct {
	name string
}

func NewSearchTableWithDatabase(ctx context.Context, db sqlite.Database) (sqlite.Table, error) {

	t, err := NewSearchTable(ctx)

	if err != nil {
		return nil, err
	}

	err = t.InitializeTable(ctx, db)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewSearchTable(ctx context.Context) (sqlite.Table, error) {

	t := SearchTable{
		name: "search",
	}

	return &t, nil
}

func (t *SearchTable) InitializeTable(ctx context.Context, db sqlite.Database) error {

	return sqlite.CreateTableIfNecessary(ctx, db, t)
}

func (t *SearchTable) Name() string {
	return t.name
}

func (t *SearchTable) Schema() string {

	schema := `CREATE VIRTUAL TABLE %s USING fts4(
		id, source, label
	);`

	// so dumb...
	return fmt.Sprintf(schema, t.Name())
}

func (t *SearchTable) IndexRecord(ctx context.Context, db sqlite.Database, i interface{}) error {
	return t.IndexRow(ctx, db, i.(map[string]string))
}

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
		return err
	}

	tx, err := conn.Begin()

	if err != nil {
		return err
	}

	s, err := tx.Prepare(fmt.Sprintf("DELETE FROM %s WHERE id = ?", t.Name()))

	if err != nil {
		return err
	}

	defer s.Close()

	_, err = s.Exec(row["id"])

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(args...)

	if err != nil {
		return err
	}

	return tx.Commit()
}
