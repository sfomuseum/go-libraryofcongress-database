package tables

import (
	"context"
	"github.com/aaronland/go-sqlite/database"
	"testing"
)

func TestIdentifiersTable(t *testing.T) {

	ctx := context.Background()

	db, err := database.NewDBWithDriver(ctx, "sqlite3", ":memory:")

	if err != nil {
		t.Fatalf("Failed to create database, %v", err)
	}

	db_t, err := NewIdentifiersTableWithDatabase(ctx, db)

	if err != nil {
		t.Fatalf("Failed to create new identifiers table, %v", err)
	}

	row := map[string]string{
		"id":     "sh00000025",
		"label":  "Women marine mammalogists",
		"source": "lcnaf",
	}

	err = db_t.IndexRecord(ctx, db, row)

	if err != nil {
		t.Fatalf("Failed to index record, %v", err)
	}
}
