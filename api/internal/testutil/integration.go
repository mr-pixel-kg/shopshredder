//go:build integration

package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenIntegrationDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open postgres connection: %v", err)
	}

	schemaName := integrationSchemaName(t)
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS " + schemaName).Error; err != nil {
		t.Fatalf("create integration schema %s: %v", schemaName, err)
	}

	scopedDB, err := gorm.Open(postgres.Open(dsn+" search_path="+schemaName), &gorm.Config{})
	if err != nil {
		t.Fatalf("open postgres connection for schema %s: %v", schemaName, err)
	}

	ApplyMigrations(t, scopedDB)

	t.Cleanup(func() {
		sqlDB, err := scopedDB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}

		if err := db.Exec("DROP SCHEMA IF EXISTS " + schemaName + " CASCADE").Error; err != nil {
			t.Fatalf("drop integration schema %s: %v", schemaName, err)
		}

		adminSQLDB, err := db.DB()
		if err == nil {
			_ = adminSQLDB.Close()
		}
	})

	return scopedDB
}

func ApplyMigrations(t *testing.T, db *gorm.DB) {
	t.Helper()

	dir := filepath.Join("..", "database", "migrations")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read migrations directory: %v", err)
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		files = append(files, filepath.Join(dir, entry.Name()))
	}
	sort.Strings(files)

	for _, path := range files {
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read migration %s: %v", path, err)
		}
		if err := db.Exec(string(sqlBytes)).Error; err != nil {
			t.Fatalf("apply migration %s: %v", path, err)
		}
	}
}

func ResetIntegrationDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	statements := []string{
		"TRUNCATE TABLE sandbox_events RESTART IDENTITY CASCADE",
		"TRUNCATE TABLE sandboxes RESTART IDENTITY CASCADE",
		"TRUNCATE TABLE images RESTART IDENTITY CASCADE",
		"TRUNCATE TABLE sessions RESTART IDENTITY CASCADE",
		"TRUNCATE TABLE audit_logs RESTART IDENTITY CASCADE",
		"TRUNCATE TABLE users RESTART IDENTITY CASCADE",
	}

	for _, stmt := range statements {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("reset database with %q: %v", stmt, err)
		}
	}
}

func integrationSchemaName(t *testing.T) string {
	t.Helper()

	// Postgres identifiers must stay simple here because we interpolate the
	// schema name into CREATE/DROP statements. A UUID suffix keeps parallel test
	// processes from colliding in CI.
	raw := fmt.Sprintf("itest_%s_%d_%s", t.Name(), time.Now().UnixNano(), uuid.NewString())
	raw = strings.ToLower(raw)

	var b strings.Builder
	for _, r := range raw {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}

	return b.String()
}
