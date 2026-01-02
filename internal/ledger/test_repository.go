package ledger

import (
	"path/filepath"
	"testing"
)

// NewTestSQLiteRepository creates a fresh SQLite DB per test
func NewTestSQLiteRepository(t *testing.T) *SQLiteRepository {
	t.Helper()

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "ledger.db")

	repo, err := NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to create test repository: %v", err)
	}

	return repo
}
