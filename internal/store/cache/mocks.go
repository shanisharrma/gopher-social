package cache

import (
	"context"

	"github.com/shanisharrma/gopher-social/internal/store"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct{}

func (m MockUserStore) Get(ctx context.Context, id int64) (*store.User, error) {
	return nil, nil
}
func (m MockUserStore) Set(ctx context.Context, user *store.User) error {
	return nil
}
