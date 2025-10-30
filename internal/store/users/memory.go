package users

import (
	"context"
)

type MemoryStorage struct {
	fakeUser *User
}

func (m *MemoryStorage) Ping(ctx context.Context) error {
	return nil
}

func (m *MemoryStorage) Get(ctx context.Context, userID int64) (*User, error) {
	return m.fakeUser, nil
}

func (m *MemoryStorage) Create(ctx context.Context) (*User, error) {
	return m.fakeUser, nil
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
