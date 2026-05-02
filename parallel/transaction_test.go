package parallel_test

import (
	"cmp"
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otakakot/sample-go-postgresql-test/database"
)

func TestTransaction_Insert(t *testing.T) {
	t.Parallel()

	dsn := cmp.Or(os.Getenv("DSN"), "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse dsn: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(t.Context(), cfg)
	if err != nil {
		t.Fatalf("failed to create pgx pool: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	beginTx, err := pool.Begin(t.Context())
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}

	t.Cleanup(func() {
		_ = beginTx.Rollback(context.Background())
	})

	tx, err := database.NewTransaction(beginTx)
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	if _, err := tx.InsertSample(t.Context(), "test"); err != nil {
		t.Fatalf("failed to insert sample: %v", err)
	}
}

func TestTransaction_Select(t *testing.T) {
	t.Parallel()

	dsn := cmp.Or(os.Getenv("DSN"), "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse dsn: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(t.Context(), cfg)
	if err != nil {
		t.Fatalf("failed to create pgx pool: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	beginTx, err := pool.Begin(t.Context())
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}

	t.Cleanup(func() {
		_ = beginTx.Rollback(context.Background())
	})

	tx, err := database.NewTransaction(beginTx)
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	want := 2

	for range want {
		if _, err := tx.InsertSample(t.Context(), "test"); err != nil {
			t.Fatalf("failed to insert sample: %v", err)
		}
	}

	samples, err := tx.ListSamples(t.Context())
	if err != nil {
		t.Fatalf("failed to list samples: %v", err)
	}

	if len(samples) != want {
		t.Fatalf("unexpected number of samples: got %d, want %d", len(samples), want)
	}
}

func TestTransaction_Delete(t *testing.T) {
	t.Parallel()

	dsn := cmp.Or(os.Getenv("DSN"), "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse dsn: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(t.Context(), cfg)
	if err != nil {
		t.Fatalf("failed to create pgx pool: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	beginTx, err := pool.Begin(t.Context())
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}

	t.Cleanup(func() {
		_ = beginTx.Rollback(context.Background())
	})

	tx, err := database.NewTransaction(beginTx)
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	for range 2 {
		if _, err := tx.InsertSample(t.Context(), "test"); err != nil {
			t.Fatalf("failed to insert sample: %v", err)
		}
	}

	if err := tx.DeleteSamples(t.Context()); err != nil {
		t.Fatalf("failed to delete samples: %v", err)
	}

	samples, err := tx.ListSamples(t.Context())
	if err != nil {
		t.Fatalf("failed to list samples: %v", err)
	}

	if len(samples) != 0 {
		t.Fatalf("expected no samples, but found %d", len(samples))
	}
}
