package database_test

import (
	"cmp"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otakakot/sample-go-postgresql-test/database"
)

func TestNewTransaction(t *testing.T) {
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

	tx, err := database.NewTransaction(pool)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

	sample, err := tx.InsertSample(t.Context(), "test")
	if err != nil {
		t.Fatalf("failed to insert sample: %v", err)
	}

	updated := "updated-test"

	if err := tx.UpdateSample(t.Context(), sample.ID, updated); err != nil {
		t.Fatalf("failed to update sample: %v", err)
	}

	found, err := tx.FindSampleByID(t.Context(), sample.ID)
	if err != nil {
		t.Fatalf("failed to find sample by ID: %v", err)
	}

	if found.ID != sample.ID || found.Name != updated {
		t.Fatalf("found sample does not match inserted sample: got %+v, want %+v", found, sample)
	}

	if samples, err := tx.ListSamples(t.Context()); err != nil {
		t.Fatalf("failed to list samples: %v", err)
	} else if len(samples) == 0 {
		t.Fatal("no samples found")
	}

	if err := tx.DeleteSamples(t.Context()); err != nil {
		t.Fatalf("failed to delete all samples: %v", err)
	}

	if samples, err := tx.ListSamples(t.Context()); err != nil {
		t.Fatalf("failed to list samples: %v", err)
	} else if len(samples) != 0 {
		t.Fatalf("expected no samples, but found %d", len(samples))
	}
}
