package domainacmes

import (
	"context"
)

type Repository interface {
	GetAll(ctx context.Context) ([]Acme, error)
	GetByID(ctx context.Context, id string) (*Acme, error)
	Create(ctx context.Context, doc *Acme) error
	Update(ctx context.Context, doc *Acme) error
	Delete(ctx context.Context, id string) error
}
