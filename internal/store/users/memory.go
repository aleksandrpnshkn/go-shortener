package users

import (
	"context"

	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

type MemoryStorage struct {
	fakeUser *User
}

func (m *MemoryStorage) Ping(ctx context.Context) error {
	return nil
}

func (m *MemoryStorage) Get(ctx context.Context, userID types.UserID) (*User, error) {
	return m.fakeUser, nil
}

func (m *MemoryStorage) Create(ctx context.Context) (types.UserID, error) {
	return m.fakeUser.ID, nil
}

func (m *MemoryStorage) Close() error {
	return nil
}

func NewMemoryStorage() *MemoryStorage {
	const fakeMemoryUserID = -1

	return &MemoryStorage{
		fakeUser: &User{
			ID: fakeMemoryUserID,
		},
	}
}
