package test_test

import (
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestTestcontainersPgx(t *testing.T) {
	migrations, err := filepath.Glob(filepath.Join("../schema", "*.sql"))
	if err != nil {
		t.Fatal(err)
	}

	container, err := postgres.Run(
		t.Context(),
		"postgres:18-alpine",
		postgres.WithInitScripts(migrations...),
		testcontainers.WithWaitStrategy(
			wait.ForAll(
				wait.ForListeningPort("5432/tcp"),
				wait.ForExec([]string{"pg_isready", "-U", "postgres", "-d", "postgres"}).
					WithPollInterval(1*time.Second).
					WithExitCodeMatcher(func(exitCode int) bool {
						return exitCode == 0
					}).
					WithStartupTimeout(30*time.Second),
			),
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	testcontainers.CleanupContainer(t, container)

	dsn, err := container.ConnectionString(t.Context(), "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	conn, err := pgxpool.ParseConfig(dsn)
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

	// if _, err := pool.Exec(t.Context(), `
	// CREATE TABLE IF NOT EXISTS samples (
	// 	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	// 	name TEXT NOT NULL,
	// 	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	// 	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
	// )`); err != nil {
	// 	t.Fatal(err)
	// }

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

func TestTestcontainersPq(t *testing.T) {
	migrations, err := filepath.Glob(filepath.Join("../schema", "*.sql"))
	if err != nil {
		t.Fatal(err)
	}

	container, err := postgres.Run(
		t.Context(),
		"postgres:18-alpine",
		postgres.WithInitScripts(migrations...),
		testcontainers.WithWaitStrategy(
			wait.ForAll(
				wait.ForListeningPort("5432/tcp"),
				wait.ForExec([]string{"pg_isready", "-U", "postgres", "-d", "postgres"}).
					WithPollInterval(1*time.Second).
					WithExitCodeMatcher(func(exitCode int) bool {
						return exitCode == 0
					}).
					WithStartupTimeout(30*time.Second),
			),
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	testcontainers.CleanupContainer(t, container)

	dsn, err := container.ConnectionString(t.Context(), "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.PingContext(t.Context()); err != nil {
		t.Fatal(err)
	}

	// if _, err := db.ExecContext(t.Context(), `
	// CREATE TABLE IF NOT EXISTS samples (
	// 	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	// 	name TEXT NOT NULL,
	// 	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	// 	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
	// )`); err != nil {
	// 	t.Fatal(err)
	// }

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
