package screening

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lcolman/fabrikam-auth-poc/pkg/definitions"
)

type SearchFilter func(definitions.Screening) bool

type ResourceOwnerIdentity string

func NewIdent(owner string) ResourceOwnerIdentity {
	return ResourceOwnerIdentity(owner)
}

type Storage interface {
	Get(context.Context, ResourceOwnerIdentity, uuid.UUID) (*definitions.Screening, error)
	Put(context.Context, ResourceOwnerIdentity, definitions.Screening) error
	Search(context.Context, ResourceOwnerIdentity, SearchFilter) ([]definitions.Screening, error)
}

// inMemory is an ephemeral in-memory implementation of the Storage persistance interface
type inMemory struct {
	data map[ResourceOwnerIdentity]map[uuid.UUID]definitions.Screening
}

func NewStorageService() Storage {
	return &inMemory{
		data: make(map[ResourceOwnerIdentity]map[uuid.UUID]definitions.Screening),
	}
}

func (i *inMemory) Get(ctx context.Context, user ResourceOwnerIdentity, screeningID uuid.UUID) (*definitions.Screening, error) {
	if i.data[user] != nil {
		result, present := i.data[user][screeningID]
		if present {
			return &result, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (i *inMemory) Put(ctx context.Context, user ResourceOwnerIdentity, item definitions.Screening) error {
	if i.data[user] == nil {
		i.data[user] = make(map[uuid.UUID]definitions.Screening)
	}
	i.data[user][*item.ID] = item
	return nil
}

// Search performs a naieve linear walk of the data filtering using the given function
func (i *inMemory) Search(ctx context.Context, user ResourceOwnerIdentity, filterFunc SearchFilter) ([]definitions.Screening, error) {
	result := []definitions.Screening{}
	if i.data[user] != nil {
		for _, item := range i.data[user] {
			if filterFunc(item) {
				result = append(result, item)
			}
		}
	}
	return result, nil
}
