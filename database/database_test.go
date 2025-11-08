package database_test

import (
	"cmp"
	"os"
	"testing"

	"github.com/otakakot/sample-go-postgresql-test/database"
)

func TestNewDatabase(t *testing.T) {
	dsn := cmp.Or(os.Getenv("DSN"), "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	pool, err := database.NewPool(dsn)
	if err != nil {
		t.Fatalf("failed to create pgx pool: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	db, err := database.NewDatabase(pool)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

	sample, err := db.InsertSample(t.Context(), "test")
	if err != nil {
		t.Fatalf("failed to insert sample: %v", err)
	}

	updated := "updated-test"

	if err := db.UpdateSample(t.Context(), sample.ID, updated); err != nil {
		t.Fatalf("failed to update sample: %v", err)
	}

	found, err := db.FindSampleByID(t.Context(), sample.ID)
	if err != nil {
		t.Fatalf("failed to find sample by ID: %v", err)
	}

	if found.ID != sample.ID || found.Name != updated {
		t.Fatalf("found sample does not match inserted sample: got %+v, want %+v", found, sample)
	}

	if samples, err := db.ListSamples(t.Context()); err != nil {
		t.Fatalf("failed to list samples: %v", err)
	} else if len(samples) == 0 {
		t.Fatal("no samples found")
	}

	if err := db.DeleteSamples(t.Context()); err != nil {
		t.Fatalf("failed to delete all samples: %v", err)
	}

	if samples, err := db.ListSamples(t.Context()); err != nil {
		t.Fatalf("failed to list samples: %v", err)
	} else if len(samples) != 0 {
		t.Fatalf("expected no samples, but found %d", len(samples))
	}
}
