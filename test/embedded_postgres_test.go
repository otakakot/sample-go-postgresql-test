package test_test

import (
	"database/sql"
	"errors"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func TestEmbeddedPostgresPgx(t *testing.T) {
	cfg := embeddedpostgres.DefaultConfig()

	database := embeddedpostgres.NewDatabase(cfg)
	if err := database.Start(); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = database.Stop()
	})

	conn, err := pgxpool.ParseConfig(cfg.GetConnectionURL() + "?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	pool, err := pgxpool.NewWithConfig(t.Context(), conn)
	if err != nil {
		t.Fatal(err)
	}

	if err := pool.Ping(t.Context()); err != nil {
		t.Fatal(err)
	}

	if _, err := pool.Exec(t.Context(), `
	CREATE TABLE IF NOT EXISTS samples (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		t.Fatal(err)
	}

	var id string

	// INSERT
	if err := pool.QueryRow(t.Context(), `INSERT INTO samples (name) VALUES ($1) RETURNING id`, "test").Scan(&id); err != nil {
		t.Fatal(err)
	}

	// SELECT
	if err := pool.QueryRow(t.Context(), `SELECT id FROM samples WHERE id = $1`, id).Scan(&id); err != nil {
		t.Fatal(err)
	}

	// UPDATE
	if _, err := pool.Exec(t.Context(), `UPDATE samples SET name = $1 WHERE id = $2`, "updated", id); err != nil {
		t.Fatal(err)
	}

	// DELETE
	if _, err := pool.Exec(t.Context(), `DELETE FROM samples WHERE id = $1`, id); err != nil {
		t.Fatal(err)
	}

	// SELECT
	if err := pool.QueryRow(t.Context(), `SELECT id FROM samples WHERE id = $1`, id).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return
		}
		t.Fatal(err)
	}
}

func TestEmbeddedPostgresPq(t *testing.T) {
	cfg := embeddedpostgres.DefaultConfig()

	database := embeddedpostgres.NewDatabase(cfg)
	if err := database.Start(); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = database.Stop()
	})

	db, err := sql.Open("postgres", cfg.GetConnectionURL()+"?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	if err := db.PingContext(t.Context()); err != nil {
		t.Fatal(err)
	}

	if _, err := db.ExecContext(t.Context(), `
	CREATE TABLE IF NOT EXISTS samples (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		t.Fatal(err)
	}

	var id string

	// INSERT
	if err := db.QueryRowContext(t.Context(), `INSERT INTO samples (name) VALUES ($1) RETURNING id`, "test").Scan(&id); err != nil {
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
