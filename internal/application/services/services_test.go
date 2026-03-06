package services

import (
	"context"
	"markitos-it-svc-goldens/internal/domain"
	"testing"
)

type fakeRepo struct{}

func (fakeRepo) GetAll(ctx context.Context) ([]domain.Golden, error) { return nil, nil }
func (fakeRepo) GetByID(ctx context.Context, id string) (*domain.Golden, error) {
	return &domain.Golden{ID: id}, nil
}
func (fakeRepo) Create(ctx context.Context, doc *domain.Golden) error { return nil }
func (fakeRepo) Update(ctx context.Context, doc *domain.Golden) error { return nil }
func (fakeRepo) Delete(ctx context.Context, id string) error          { return nil }

func FuzzNewGoldenService(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) < 1 {
			return
		}
		_ = NewGoldenService(nil)
	})
}

func FuzzGoldenService_GetAllGoldens(f *testing.F) {
	f.Add([]byte("seed"))
	f.Fuzz(func(t *testing.T, _ []byte) {
		s := NewGoldenService(fakeRepo{})
		_, err := s.GetAllGoldens(context.Background())
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
	})
}

func FuzzGoldenService_GetGoldenByID(f *testing.F) {
	f.Add("test-id")
	f.Fuzz(func(t *testing.T, id string) {
		s := NewGoldenService(fakeRepo{})
		got, err := s.GetGoldenByID(context.Background(), id)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if got == nil || got.ID != id {
			t.Fatalf("unexpected result: %+v", got)
		}
	})
}

func FuzzGoldenService_CreateGolden(f *testing.F) {
	f.Add([]byte("seed"))
	f.Fuzz(func(t *testing.T, _ []byte) {
		s := NewGoldenService(fakeRepo{})
		err := s.CreateGolden(context.Background(), &domain.Golden{})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
	})
}

func FuzzGoldenService_UpdateGolden(f *testing.F) {
	f.Add([]byte("seed"))
	f.Fuzz(func(t *testing.T, _ []byte) {
		s := NewGoldenService(fakeRepo{})
		err := s.UpdateGolden(context.Background(), &domain.Golden{})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
	})
}

func FuzzGoldenService_DeleteGolden(f *testing.F) {
	f.Add("test-id")
	f.Fuzz(func(t *testing.T, id string) {
		s := NewGoldenService(fakeRepo{})
		err := s.DeleteGolden(context.Background(), id)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
	})
}
