package postgres

import (
	"context"
	"database/sql"
	"markitos-it-svc-goldens/internal/domain"
	"reflect"
	"testing"
	"time"
)

func helperClosedDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("postgres", "host=127.0.0.1 port=1 user=test dbname=test sslmode=disable connect_timeout=1")
	if err != nil {
		t.Fatalf("failed to create db handle: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("failed to close db handle: %v", err)
	}

	return db
}

func helperRandomGolden(t *testing.T) *domain.Golden {
	t.Helper()

	prefix := domain.HelperRandomAlphaPrefix(t, 8)
	return &domain.Golden{
		ID:          prefix + "-golden-id",
		Title:       prefix + "-golden-title",
		Description: prefix + "-golden-description",
		Category:    prefix + "-golden-category",
		Tags:        []string{prefix + "-go", prefix + "-grpc", prefix + "-postgres"},
		UpdatedAt:   time.Date(2026, 3, 6, 12, 0, 0, 0, time.UTC),
		ContentB64:  prefix + "-Y29udGVudA==",
		CoverImage:  "https://example.com/" + prefix + "/cover.png",
	}
}

func TestNewGoldenRepository(t *testing.T) {
	prefix := domain.HelperRandomAlphaPrefix(t, 6)
	db := helperClosedDB(t)

	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name   string
		args   args
		wantDB *sql.DB
	}{
		{
			name:   prefix + "-build-repository-with-same-db",
			args:   args{db: db},
			wantDB: db,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGoldenRepository(tt.args.db)
			if got == nil {
				t.Fatalf("NewGoldenRepository() returned nil")
			}
			if got.db != tt.wantDB {
				t.Errorf("NewGoldenRepository().db = %v, want %v", got.db, tt.wantDB)
			}
		})
	}
}

func TestGoldenRepository_InitSchema(t *testing.T) {
	prefix := domain.HelperRandomAlphaPrefix(t, 6)
	db := helperClosedDB(t)

	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    prefix + "-returns-error-on-closed-db",
			fields:  fields{db: db},
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &GoldenRepository{db: tt.fields.db}
			if err := r.InitSchema(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("GoldenRepository.InitSchema() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGoldenRepository_SeedData(t *testing.T) {
	prefix := domain.HelperRandomAlphaPrefix(t, 6)
	db := helperClosedDB(t)

	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    prefix + "-returns-error-on-closed-db",
			fields:  fields{db: db},
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &GoldenRepository{db: tt.fields.db}
			if err := r.SeedData(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("GoldenRepository.SeedData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGoldenRepository_GetAll(t *testing.T) {
	prefix := domain.HelperRandomAlphaPrefix(t, 6)
	db := helperClosedDB(t)

	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Golden
		wantErr bool
	}{
		{
			name:    prefix + "-returns-error-on-closed-db",
			fields:  fields{db: db},
			args:    args{ctx: context.Background()},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &GoldenRepository{db: tt.fields.db}
			got, err := r.GetAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GoldenRepository.GetAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoldenRepository.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoldenRepository_GetByID(t *testing.T) {
	prefix := domain.HelperRandomAlphaPrefix(t, 6)
	db := helperClosedDB(t)

	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Golden
		wantErr bool
	}{
		{
			name:    prefix + "-returns-error-on-closed-db",
			fields:  fields{db: db},
			args:    args{ctx: context.Background(), id: prefix + "-missing-id"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &GoldenRepository{db: tt.fields.db}
			got, err := r.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GoldenRepository.GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoldenRepository.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoldenRepository_Create(t *testing.T) {
	prefix := domain.HelperRandomAlphaPrefix(t, 6)
	db := helperClosedDB(t)
	randomDoc := helperRandomGolden(t)

	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx context.Context
		doc *domain.Golden
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    prefix + "-returns-error-on-closed-db",
			fields:  fields{db: db},
			args:    args{ctx: context.Background(), doc: randomDoc},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &GoldenRepository{db: tt.fields.db}
			if err := r.Create(tt.args.ctx, tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("GoldenRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGoldenRepository_Update(t *testing.T) {
	prefix := domain.HelperRandomAlphaPrefix(t, 6)
	db := helperClosedDB(t)
	randomDoc := helperRandomGolden(t)

	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx context.Context
		doc *domain.Golden
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    prefix + "-returns-error-on-closed-db",
			fields:  fields{db: db},
			args:    args{ctx: context.Background(), doc: randomDoc},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &GoldenRepository{db: tt.fields.db}
			if err := r.Update(tt.args.ctx, tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("GoldenRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGoldenRepository_Delete(t *testing.T) {
	prefix := domain.HelperRandomAlphaPrefix(t, 6)
	db := helperClosedDB(t)

	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    prefix + "-returns-error-on-closed-db",
			fields:  fields{db: db},
			args:    args{ctx: context.Background(), id: prefix + "-to-delete"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &GoldenRepository{db: tt.fields.db}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("GoldenRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func FuzzGoldenRepository_GetByID_ClosedDB(f *testing.F) {
	f.Add("missing-id")
	f.Add("")
	f.Add("id-con-123")

	f.Fuzz(func(t *testing.T, id string) {
		db := helperClosedDB(t)
		r := &GoldenRepository{db: db}

		_, err := r.GetByID(context.Background(), id)
		if err == nil {
			t.Fatalf("expected error for closed db, got nil (id=%q)", id)
		}
	})
}

func FuzzGoldenRepository_Delete_ClosedDB(f *testing.F) {
	f.Add("to-delete")
	f.Add("")
	f.Add("abc-123")

	f.Fuzz(func(t *testing.T, id string) {
		db := helperClosedDB(t)
		r := &GoldenRepository{db: db}

		err := r.Delete(context.Background(), id)
		if err == nil {
			t.Fatalf("expected error for closed db, got nil (id=%q)", id)
		}
	})
}

func FuzzGoldenRepository_Create_ClosedDB(f *testing.F) {
	f.Add("id-1", "title", "desc", "cat", "Y29udGVudA==", "https://example.com/cover.png")
	f.Add("", "", "", "", "", "")

	f.Fuzz(func(t *testing.T, id, title, description, category, contentB64, coverImage string) {
		db := helperClosedDB(t)
		r := &GoldenRepository{db: db}

		doc := &domain.Golden{
			ID:          id,
			Title:       title,
			Description: description,
			Category:    category,
			Tags:        []string{"a", "b"},
			UpdatedAt:   time.Now(),
			ContentB64:  contentB64,
			CoverImage:  coverImage,
		}

		err := r.Create(context.Background(), doc)
		if err == nil {
			t.Fatalf("expected error for closed db, got nil")
		}
	})
}
