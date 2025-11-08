package parallel_test

import (
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/otakakot/sample-go-postgresql-test/database"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestCreatePostgreSQL_Insert(t *testing.T) {
	t.Parallel()

	_, filename, _, _ := runtime.Caller(0)
	parallelDir := filepath.Dir(filename)
	rootDir := filepath.Dir(parallelDir)
	schemaDir := filepath.Join(rootDir, "schema")

	dsn := CreatePostgreSQL(t, schemaDir)

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse dsn: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(t.Context(), cfg)
	if err != nil {
		t.Fatalf("failed to create pgx pool: %v", err)
	}

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

func TestCreatePostgreSQL_Select(t *testing.T) {
	t.Parallel()

	_, filename, _, _ := runtime.Caller(0)
	parallelDir := filepath.Dir(filename)
	rootDir := filepath.Dir(parallelDir)
	schemaDir := filepath.Join(rootDir, "schema")

	dsn := CreatePostgreSQL(t, schemaDir)

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse dsn: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(t.Context(), cfg)
	if err != nil {
		t.Fatalf("failed to create pgx pool: %v", err)
	}

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

func TestCreatePostgreSQL_Delete(t *testing.T) {
	t.Parallel()

	_, filename, _, _ := runtime.Caller(0)
	parallelDir := filepath.Dir(filename)
	rootDir := filepath.Dir(parallelDir)
	schemaDir := filepath.Join(rootDir, "schema")

	dsn := CreatePostgreSQL(t, schemaDir)

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse dsn: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(t.Context(), cfg)
	if err != nil {
		t.Fatalf("failed to create pgx pool: %v", err)
	}

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

func CreatePostgreSQL(
	t *testing.T,
	sqlDir string,
) string {
	t.Helper()

	migrations, err := filepath.Glob(filepath.Join(sqlDir, "*.sql"))
	if err != nil {
		t.Fatal(err)
	}

	if len(migrations) == 0 {
		t.Fatalf("No SQL files found in schema directory: %s", sqlDir)
	}

	container, err := postgres.Run(
		t.Context(),
		"postgres:18-alpine",
		postgres.WithInitScripts(migrations...),
		testcontainers.WithEnv(map[string]string{
			"TZ":                        "UTC",
			"LANG":                      "ja_JP.UTF-8",
			"POSTGRES_INITDB_ARGS":      "--encoding=UTF-8",
			"POSTGRES_HOST_AUTH_METHOD": "trust",
		}),
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

	return dsn
}
