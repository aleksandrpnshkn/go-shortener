package users

import (
	"context"

	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

// MemoryStorage - in memory хранилище пользователей.
type MemoryStorage struct {
	fakeUser *User
}

// Ping заглушка для проверки доступности хранилища.
func (m *MemoryStorage) Ping(ctx context.Context) error {
	return nil
}

// Get "достаёт" фейкового пользователя из хранилища.
func (m *MemoryStorage) Get(ctx context.Context, userID types.UserID) (*User, error) {
	return m.fakeUser, nil
}

// Create "создаёт" фейкового пользователя в хранилища.
func (m *MemoryStorage) Create(ctx context.Context) (types.UserID, error) {
	return m.fakeUser.ID, nil
}

// Close заглушка для закрытия хранилища.
func (m *MemoryStorage) Close() error {
	return nil
}

// NewMemoryStorage создаёт новое in memory хранилище пользователей.
func NewMemoryStorage() *MemoryStorage {
	const fakeMemoryUserID = -1

	return &MemoryStorage{
		fakeUser: &User{
			ID: fakeMemoryUserID,
		},
	}
}
