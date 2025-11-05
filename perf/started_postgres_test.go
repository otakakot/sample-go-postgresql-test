package perf_test

import (
	"cmp"
	"database/sql"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func TestStartedPostgresPgx(t *testing.T) {
	dsn := cmp.Or(os.Getenv("DSN"), "postgres://postgres:password@localhost:5432/postgres?sslmode=disable")

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
}

func TestStartedPostgresPq(t *testing.T) {
	dsn := cmp.Or(os.Getenv("DSN"), "postgres://postgres:password@localhost:5432/postgres?sslmode=disable")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.PingContext(t.Context()); err != nil {
		t.Fatal(err)
	}
}
