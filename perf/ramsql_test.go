package perf_test

import (
	"database/sql"
	"testing"

	_ "github.com/proullon/ramsql/driver"
)

func TestRamsql(t *testing.T) {
	db, err := sql.Open("ramsql", "TestRamsql")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	if err := db.PingContext(t.Context()); err != nil {
		t.Fatal(err)
	}
}
