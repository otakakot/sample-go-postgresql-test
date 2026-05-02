package test_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	_ "github.com/proullon/ramsql/driver"
)

func TestRamsql(t *testing.T) {
	db, err := sql.Open("ramsql", "TestRamsql")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	if err := db.PingContext(t.Context()); err != nil {
		t.Fatal(err)
	}

	// Not Supported: gen_random_uuid(), CURRENT_TIMESTAMP
	if _, err = db.ExecContext(t.Context(), `
	CREATE TABLE IF NOT EXISTS samples (
    	id UUID PRIMARY KEY,
    	name TEXT NOT NULL,
    	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	)`); err != nil {
		t.Fatal(err)
	}

	id := uuid.NewString()

	// INSERT
	if _, err := db.ExecContext(t.Context(), `INSERT INTO samples (id, name) VALUES ($1, $2) RETURNING id`, id, "test"); err != nil {
		t.Fatal(err)
	}

	// SELECT
	if err := db.QueryRowContext(t.Context(), `SELECT id FROM samples WHERE id = $1`, id).Scan(&id); err != nil {
		t.Fatal(err)
	}

	// UPDATE
	if _, err := db.ExecContext(t.Context(), `UPDATE samples SET name = $1 WHERE id = $2`, "updated", id); err != nil {
		t.Fatal(err)
	}

	// DELETE
	if _, err := db.ExecContext(t.Context(), `DELETE FROM samples WHERE id = $1`, id); err != nil {
		t.Fatal(err)
	}

	// SELECT
	if err := db.QueryRowContext(t.Context(), `SELECT id FROM samples WHERE id = $1`, id).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
		t.Fatal(err)
	}
}
