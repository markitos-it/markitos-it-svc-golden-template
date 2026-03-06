package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"markitos-it-svc-goldens/internal/domain"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

func helperIntegrationDB(t *testing.T) *sql.DB {
	t.Helper()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")
	if host == "" || port == "" || user == "" || name == "" {
		t.Skip("integration DB not configured; set DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_NAME")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, name)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to open integration db: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	if err := db.Ping(); err != nil {
		t.Skipf("integration DB unreachable: %v", err)
	}

	return db
}

func helperEnsureSchemaAndClean(t *testing.T, r *GoldenRepository) {
	t.Helper()
	ctx := context.Background()
	if err := r.InitSchema(ctx); err != nil {
		t.Fatalf("InitSchema() failed: %v", err)
	}
	if _, err := r.db.ExecContext(ctx, "TRUNCATE TABLE goldens"); err != nil {
		t.Fatalf("failed to truncate goldens: %v", err)
	}
}

func helperInsertDocDirect(t *testing.T, db *sql.DB, doc *domain.Golden) {
	t.Helper()
	_, err := db.Exec(
		`INSERT INTO goldens (id, title, description, category, tags, updated_at, content_b64, cover_image)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		doc.ID, doc.Title, doc.Description, doc.Category, pqStringArray(doc.Tags), doc.UpdatedAt, doc.ContentB64, doc.CoverImage,
	)
	if err != nil {
		t.Fatalf("failed to insert test row: %v", err)
	}
}

func pqStringArray(tags []string) interface{} {
	if len(tags) == 0 {
		return "{}"
	}
	quoted := make([]string, 0, len(tags))
	for _, tag := range tags {
		escaped := strings.ReplaceAll(tag, "\"", "\\\"")
		quoted = append(quoted, fmt.Sprintf("\"%s\"", escaped))
	}
	return "{" + strings.Join(quoted, ",") + "}"
}

func TestGoldenRepository_InitSchema_Success_Integration(t *testing.T) {
	db := helperIntegrationDB(t)
	r := NewGoldenRepository(db)

	if err := r.InitSchema(context.Background()); err != nil {
		t.Fatalf("InitSchema() unexpected error: %v", err)
	}
}

func TestGoldenRepository_SeedData_CountGreaterThanZero_Integration(t *testing.T) {
	db := helperIntegrationDB(t)
	r := NewGoldenRepository(db)
	helperEnsureSchemaAndClean(t, r)

	doc := helperRandomGolden(t)
	helperInsertDocDirect(t, db, doc)

	if err := r.SeedData(context.Background()); err != nil {
		t.Fatalf("SeedData() unexpected error: %v", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM goldens").Scan(&count); err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected count=1 when pre-seeded, got %d", count)
	}
}

func TestGoldenRepository_SeedData_Success_Integration(t *testing.T) {
	db := helperIntegrationDB(t)
	r := NewGoldenRepository(db)
	helperEnsureSchemaAndClean(t, r)

	if err := r.SeedData(context.Background()); err != nil {
		t.Fatalf("SeedData() unexpected error: %v", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM goldens").Scan(&count); err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected seeded count=2, got %d", count)
	}
}

func TestGoldenRepository_GetAll_Success_Integration(t *testing.T) {
	db := helperIntegrationDB(t)
	r := NewGoldenRepository(db)
	helperEnsureSchemaAndClean(t, r)

	doc := helperRandomGolden(t)
	helperInsertDocDirect(t, db, doc)

	docs, err := r.GetAll(context.Background())
	if err != nil {
		t.Fatalf("GetAll() unexpected error: %v", err)
	}
	if len(docs) != 1 {
		t.Fatalf("expected 1 doc, got %d", len(docs))
	}
}

func TestGoldenRepository_GetByID_NotFound_Integration(t *testing.T) {
	db := helperIntegrationDB(t)
	r := NewGoldenRepository(db)
	helperEnsureSchemaAndClean(t, r)

	_, err := r.GetByID(context.Background(), "missing-id")
	if err == nil {
		t.Fatalf("GetByID() expected not found error")
	}
}

func TestGoldenRepository_GetByID_Success_Integration(t *testing.T) {
	db := helperIntegrationDB(t)
	r := NewGoldenRepository(db)
	helperEnsureSchemaAndClean(t, r)

	doc := helperRandomGolden(t)
	helperInsertDocDirect(t, db, doc)

	got, err := r.GetByID(context.Background(), doc.ID)
	if err != nil {
		t.Fatalf("GetByID() unexpected error: %v", err)
	}
	if got == nil || got.ID != doc.ID {
		t.Fatalf("expected id=%s, got %+v", doc.ID, got)
	}
}

func TestGoldenRepository_Create_Success_Integration(t *testing.T) {
	db := helperIntegrationDB(t)
	r := NewGoldenRepository(db)
	helperEnsureSchemaAndClean(t, r)

	doc := helperRandomGolden(t)
	if err := r.Create(context.Background(), doc); err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM goldens WHERE id = $1", doc.ID).Scan(&count); err != nil {
		t.Fatalf("verification query failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 row for id=%s, got %d", doc.ID, count)
	}
}

func TestGoldenRepository_Update_NotFound_Integration(t *testing.T) {
	db := helperIntegrationDB(t)
	r := NewGoldenRepository(db)
	helperEnsureSchemaAndClean(t, r)

	doc := helperRandomGolden(t)
	if err := r.Update(context.Background(), doc); err == nil {
		t.Fatalf("Update() expected not found error")
	}
}

func TestGoldenRepository_Update_Success_Integration(t *testing.T) {
	db := helperIntegrationDB(t)
	r := NewGoldenRepository(db)
	helperEnsureSchemaAndClean(t, r)

	doc := helperRandomGolden(t)
	helperInsertDocDirect(t, db, doc)

	doc.Title = doc.Title + "-updated"
	if err := r.Update(context.Background(), doc); err != nil {
		t.Fatalf("Update() unexpected error: %v", err)
	}

	var title string
	if err := db.QueryRow("SELECT title FROM goldens WHERE id = $1", doc.ID).Scan(&title); err != nil {
		t.Fatalf("verification query failed: %v", err)
	}
	if title != doc.Title {
		t.Fatalf("expected updated title=%s, got %s", doc.Title, title)
	}
}

func TestGoldenRepository_Delete_NotFound_Integration(t *testing.T) {
	db := helperIntegrationDB(t)
	r := NewGoldenRepository(db)
	helperEnsureSchemaAndClean(t, r)

	if err := r.Delete(context.Background(), "missing-id"); err == nil {
		t.Fatalf("Delete() expected not found error")
	}
}

func TestGoldenRepository_Delete_Success_Integration(t *testing.T) {
	db := helperIntegrationDB(t)
	r := NewGoldenRepository(db)
	helperEnsureSchemaAndClean(t, r)

	doc := helperRandomGolden(t)
	helperInsertDocDirect(t, db, doc)

	if err := r.Delete(context.Background(), doc.ID); err != nil {
		t.Fatalf("Delete() unexpected error: %v", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM goldens WHERE id = $1", doc.ID).Scan(&count); err != nil {
		t.Fatalf("verification query failed: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 rows for id=%s, got %d", doc.ID, count)
	}
}
