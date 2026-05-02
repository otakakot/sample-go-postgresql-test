package perf_test

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestTestcontainersPgx(t *testing.T) {
	container, err := postgres.Run(
		t.Context(),
		"postgres:18-alpine",
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

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatal(err)
	}

	pool, err := pgxpool.NewWithConfig(t.Context(), cfg)
	if err != nil {
		t.Fatal(err)
	}

	if err := pool.Ping(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func TestTestcontainersPq(t *testing.T) {
	container, err := postgres.Run(
		t.Context(),
		"postgres:18-alpine",
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

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatal(err)
	}

	pool, err := pgxpool.NewWithConfig(t.Context(), cfg)
	if err != nil {
		t.Fatal(err)
	}

	if err := pool.Ping(t.Context()); err != nil {
		t.Fatal(err)
	}
}
