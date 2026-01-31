package database

import (
	"testing"

	cfg "github.com/artrctx/quoin-core/internal/config"
	"github.com/artrctx/quoin-core/tests/dbtest"
)

var testDBConfig *cfg.DatabaseConfig = &cfg.DatabaseConfig{}

func TestMain(m *testing.M) {
	dbtest.RunWithPostgres(m, testDBConfig)
}

func TestServiceNew(t *testing.T) {
	srv := New(testDBConfig.ConnStr())
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestServiceHealth(t *testing.T) {
	srv := New(testDBConfig.ConnStr())

	stats := srv.Health()

	if stats["status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["status"])
	}

	if _, ok := stats["error"]; ok {
		t.Fatalf("expected error not to be present")
	}

	if stats["message"] != "The database is healthy." {
		t.Fatalf("expected message to be 'The database is healthy.', got %s", stats["message"])
	}
}

func TestServiceClose(t *testing.T) {
	srv := New(testDBConfig.ConnStr())
	if srv.Close() != nil {
		t.Fatal("expected Service Close() to returned nil")
	}
}
