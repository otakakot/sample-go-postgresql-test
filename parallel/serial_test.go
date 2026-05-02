package parallel_test

import (
	"cmp"
	"context"
	"os"
	"testing"

	"github.com/otakakot/sample-go-postgresql-test/database"
)

func TestSerial_Insert(t *testing.T) {
	dsn := cmp.Or(os.Getenv("DSN"), "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	pool, err := database.NewPool(dsn)
	if err != nil {
		t.Fatalf("failed to create pgx pool: %v", err)
	}

	db, err := database.NewDatabase(pool)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `TRUNCATE TABLE samples`)

		pool.Close()
	})

	if _, err := db.InsertSample(t.Context(), "test"); err != nil {
		t.Fatalf("failed to insert sample: %v", err)
	}

	samples, err := db.ListSamples(t.Context())
	if err != nil {
		t.Fatalf("failed to list samples: %v", err)
	}

	if len(samples) != 1 {
		t.Fatalf("unexpected number of samples: got %d, want %d", len(samples), 1)
	}
}

func TestSerial_Select(t *testing.T) {
	dsn := cmp.Or(os.Getenv("DSN"), "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	pool, err := database.NewPool(dsn)
	if err != nil {
		t.Fatalf("failed to create pgx pool: %v", err)
	}

	db, err := database.NewDatabase(pool)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `TRUNCATE TABLE samples`)

		pool.Close()
	})

	want := 2

	for range want {
		if _, err := db.InsertSample(t.Context(), "test"); err != nil {
			t.Fatalf("failed to insert sample: %v", err)
		}
	}

	samples, err := db.ListSamples(t.Context())
	if err != nil {
		t.Fatalf("failed to list samples: %v", err)
	}

	if len(samples) != want {
		t.Fatalf("unexpected number of samples: got %d, want %d", len(samples), want)
	}
}

func TestSerial_Delete(t *testing.T) {
	dsn := cmp.Or(os.Getenv("DSN"), "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	pool, err := database.NewPool(dsn)
	if err != nil {
		t.Fatalf("failed to create pgx pool: %v", err)
	}

	db, err := database.NewDatabase(pool)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `TRUNCATE TABLE samples`)

		pool.Close()
	})

	for range 2 {
		if _, err := db.InsertSample(t.Context(), "test"); err != nil {
			t.Fatalf("failed to insert sample: %v", err)
		}
	}

	if err := db.DeleteSamples(t.Context()); err != nil {
		t.Fatalf("failed to delete samples: %v", err)
	}

	samples, err := db.ListSamples(t.Context())
	if err != nil {
		t.Fatalf("failed to list samples: %v", err)
	}

	if len(samples) != 0 {
		t.Fatalf("expected no samples, but found %d", len(samples))
	}
}
