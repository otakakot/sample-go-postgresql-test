package perf_test

import (
	"database/sql"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
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
}
