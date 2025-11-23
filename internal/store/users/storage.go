// Package users это интерфейс для работы с хранилищем пользователей.
package users

import (
	"context"
	"errors"

	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

// GuestID - null-value для работы с id пользователя.
// Означает что пользователя не удалось аутентифицировать.
const GuestID = types.UserID(0)

// User - сущность юзера в хранилище.
type User struct {
	ID types.UserID
}

// Storage - интерфейс хранилища пользователей.
type Storage interface {
	// Ping проверяет доступность хранилища.
	Ping(ctx context.Context) error

	// Get достаёт пользователя из хранилища.
	Get(ctx context.Context, userID types.UserID) (*User, error)

	// Create создаёт нового пользователя в хранилище.
	Create(ctx context.Context) (types.UserID, error)

	// Close закрывает соединение с хранилищем, если этого требует реализация.
	// Вызывается при завершении работы программы.
	Close() error
}

// Ошибки работы с хранилищем пользователей.
var (
	ErrUserNotFound = errors.New("user not found")
)
