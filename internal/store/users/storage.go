package users

import (
	"context"
	"errors"

	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

// GuestID - null-value для работы с id пользователя.
// Означает что пользователя не удалось аутентифицировать.
const GuestID = types.UserID(0)

type User struct {
	ID types.UserID
}

type Storage interface {
	// Ping проверяет доступность БД.
	Ping(ctx context.Context) error

	// Get достаёт пользователя из БД.
	Get(ctx context.Context, userID types.UserID) (*User, error)

	// Create создаёт нового пользователя в БД.
	Create(ctx context.Context) (types.UserID, error)

	// Close закрывает соединение с хранилищем, если этого требует реализация.
	// Вызывается при завершении работы программы.
	Close() error
}

var (
	ErrUserNotFound = errors.New("user not found")
)
