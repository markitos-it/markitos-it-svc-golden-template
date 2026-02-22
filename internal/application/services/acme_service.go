package services

import (
	"context"
	domaingoldens "markitos-it-svc-goldens/internal/domain/domainacmes"
)

type GoldenService struct {
	repo domaingoldens.Repository
}

func NewGoldenService(repo domaingoldens.Repository) *GoldenService {
	return &GoldenService{
		repo: repo,
	}
}

func (s *GoldenService) GetAllGoldens(ctx context.Context) ([]domaingoldens.Golden, error) {
	return s.repo.GetAll(ctx)
}

func (s *GoldenService) GetGoldenByID(ctx context.Context, id string) (*domaingoldens.Golden, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *GoldenService) CreateGolden(ctx context.Context, doc *domaingoldens.Golden) error {
	return s.repo.Create(ctx, doc)
}

func (s *GoldenService) UpdateGolden(ctx context.Context, doc *domaingoldens.Golden) error {
	return s.repo.Update(ctx, doc)
}

func (s *GoldenService) DeleteGolden(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
