package parallel_test

import (
	"cmp"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/otakakot/sample-go-postgresql-test/database"
)

func TestCreateDatabase_Insert(t *testing.T) {
	t.Parallel()

	host := cmp.Or(os.Getenv("POSTGRES_HOST"), "localhost")
	port := cmp.Or(os.Getenv("POSTGRES_PORT"), "5432")
	user := cmp.Or(os.Getenv("POSTGRES_USER"), "postgres")
	password := cmp.Or(os.Getenv("POSTGRES_PASSWORD"), "postgres")
	dbName := cmp.Or(os.Getenv("POSTGRES_DB_NAME"), "postgres")

	dsn := CreateDatabase(t, host, port, user, password, dbName)

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

	_, filename, _, _ := runtime.Caller(0)
	parallelDir := filepath.Dir(filename)
	rootDir := filepath.Dir(parallelDir)
	schemaDir := filepath.Join(rootDir, "schema")

	Migrate(t, pool, schemaDir)

	db, err := database.NewDatabase(pool)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

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

func TestCreateDatabase_Select(t *testing.T) {
	t.Parallel()

	host := cmp.Or(os.Getenv("POSTGRES_HOST"), "localhost")
	port := cmp.Or(os.Getenv("POSTGRES_PORT"), "5432")
	user := cmp.Or(os.Getenv("POSTGRES_USER"), "postgres")
	password := cmp.Or(os.Getenv("POSTGRES_PASSWORD"), "postgres")
	dbName := cmp.Or(os.Getenv("POSTGRES_DB_NAME"), "postgres")

	dsn := CreateDatabase(t, host, port, user, password, dbName)

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

	_, filename, _, _ := runtime.Caller(0)
	parallelDir := filepath.Dir(filename)
	rootDir := filepath.Dir(parallelDir)
	schemaDir := filepath.Join(rootDir, "schema")

	Migrate(t, pool, schemaDir)

	db, err := database.NewDatabase(pool)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

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

func TestCreateDatabase_Delete(t *testing.T) {
	t.Parallel()

	host := cmp.Or(os.Getenv("POSTGRES_HOST"), "localhost")
	port := cmp.Or(os.Getenv("POSTGRES_PORT"), "5432")
	user := cmp.Or(os.Getenv("POSTGRES_USER"), "postgres")
	password := cmp.Or(os.Getenv("POSTGRES_PASSWORD"), "postgres")
	dbName := cmp.Or(os.Getenv("POSTGRES_DB_NAME"), "postgres")

	dsn := CreateDatabase(t, host, port, user, password, dbName)

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

	_, filename, _, _ := runtime.Caller(0)
	parallelDir := filepath.Dir(filename)
	rootDir := filepath.Dir(parallelDir)
	schemaDir := filepath.Join(rootDir, "schema")

	Migrate(t, pool, schemaDir)

	db, err := database.NewDatabase(pool)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

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

func CreateDatabase(
	t *testing.T,
	host string,
	port string,
	user string,
	password string,
	db string,
) string {
	t.Helper()

	name := "test_" + strings.ReplaceAll(uuid.NewString(), "-", "_")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, db)

	conn, err := pgx.Connect(t.Context(), dsn)
	if err != nil {
		t.Fatalf("failed to connect to postgres database: %v", err)
	}

	t.Cleanup(func() {
		_ = conn.Close(context.Background())
	})

	if _, err := conn.Exec(t.Context(), "CREATE DATABASE "+name+" WITH OWNER = "+user+" ENCODING = 'UTF8'"); err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

	t.Cleanup(func() {
		_, _ = conn.Exec(context.Background(), "DROP DATABASE IF EXISTS "+name)
	})

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, name)
}

func Migrate(
	t *testing.T,
	pool *pgxpool.Pool,
	sqlDir string,
) {
	t.Helper()

	migrations, err := filepath.Glob(filepath.Join(sqlDir, "*.sql"))
	if err != nil {
		t.Fatal(err)
	}

	if len(migrations) == 0 {
		t.Fatalf("No SQL files found in schema directory: %s", sqlDir)
	}

	for _, migration := range migrations {
		sqlContent, err := os.ReadFile(migration)
		if err != nil {
			t.Fatalf("failed to read migration file %s: %v", migration, err)
		}

		if _, err := pool.Exec(t.Context(), string(sqlContent)); err != nil {
			t.Fatalf("failed to apply migration %s: %v", migration, err)
		}
	}
}
