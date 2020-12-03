package screening

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lcolman/fabrikam-auth-poc/pkg/definitions"
)

type Service interface {
	Find(ctx context.Context, client ResourceOwnerIdentity, id *uuid.UUID) (definitions.Screening, error)
	List(ctx context.Context, client ResourceOwnerIdentity) ([]definitions.Screening, error)
	Create(context.Context, ResourceOwnerIdentity, definitions.Screening) (*uuid.UUID, error)
	Update(context.Context, ResourceOwnerIdentity, definitions.Screening) error
	Delete(ctx context.Context, client ResourceOwnerIdentity, id *uuid.UUID) error
}

type service struct {
	storage Storage
}

func NewScreeningService(storageService Storage) Service {
	return &service{
		storage: storageService,
	}
}

func (s *service) Find(ctx context.Context, client ResourceOwnerIdentity, id *uuid.UUID) (definitions.Screening, error) {
	if id == nil {
		return definitions.Screening{}, fmt.Errorf("invalid screening ID")
	}
	item, err := s.storage.Get(ctx, client, *id)
	if err != nil {
		return definitions.Screening{}, err
	}
	return *item, nil
}

func (s *service) List(ctx context.Context, client ResourceOwnerIdentity) ([]definitions.Screening, error) {
	showAll := func(_ definitions.Screening) bool {
		return true
	}
	return s.storage.Search(ctx, client, showAll)
}

func (s *service) Create(ctx context.Context, client ResourceOwnerIdentity, new definitions.Screening) (*uuid.UUID, error) {
	if new.ID == nil {
		newUUID := uuid.New()
		new.ID = &newUUID
	}
	_, err := s.Find(ctx, client, new.ID)
	if err == nil {
		return nil, fmt.Errorf("already exists")
	}
	if err.Error() != "not found" {
		return nil, err
	}
	return new.ID, s.storage.Put(ctx, client, new)
}

func (s *service) Update(ctx context.Context, client ResourceOwnerIdentity, new definitions.Screening) error {
	return fmt.Errorf("Not Yet Implemented")
}

func (s *service) Delete(ctx context.Context, client ResourceOwnerIdentity, id *uuid.UUID) error {
	return fmt.Errorf("Not Yet Implemented")
}
