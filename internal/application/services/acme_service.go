package services

import (
	"context"
	"markitos-it-svc-acmes/internal/domain/domainacmes"
)

type AcmeService struct {
	repo domainacmes.Repository
}

func NewAcmeService(repo domainacmes.Repository) *AcmeService {
	return &AcmeService{
		repo: repo,
	}
}

func (s *AcmeService) GetAllAcmes(ctx context.Context) ([]domainacmes.Acme, error) {
	return s.repo.GetAll(ctx)
}

func (s *AcmeService) GetAcmeByID(ctx context.Context, id string) (*domainacmes.Acme, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AcmeService) CreateAcme(ctx context.Context, doc *domainacmes.Acme) error {
	return s.repo.Create(ctx, doc)
}

func (s *AcmeService) UpdateAcme(ctx context.Context, doc *domainacmes.Acme) error {
	return s.repo.Update(ctx, doc)
}

func (s *AcmeService) DeleteAcme(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
