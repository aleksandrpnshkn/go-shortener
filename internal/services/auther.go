package services

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v4"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

// Claims - информация, зашитая в токен.
type Claims struct {
	jwt.RegisteredClaims
	UserID int64
}

// Auther - интерфейс для аутентификации пользователей.
type Auther interface {
	// ParseToken парсит токен.
	ParseToken(ctx context.Context, token string) (types.UserID, error)

	// RegisterUser регистрирует юзера.
	RegisterUser(ctx context.Context) (types.UserID, string, error)

	// FromUserContext достаёт юзера контекста.
	FromUserContext(ctx context.Context) (types.UserID, error)
}

// JwtAuther - реализация для аутентификации пользователей на основе JWT-токена.
type JwtAuther struct {
	usersStorage users.Storage

	secretKey string
}

type ctxKey string

// Ошибки JwtAuther.
var (
	ErrInvalidToken = errors.New("invalid token")
)

const ctxUserID ctxKey = "user"

// ParseToken парсит JWT-токен.
func (a *JwtAuther) ParseToken(ctx context.Context, tokenString string) (types.UserID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(a.secretKey), nil
	})
	if err != nil {
		return users.GuestID, ErrInvalidToken
	}

	claims := token.Claims.(*Claims)

	return types.UserID(claims.UserID), nil
}

// RegisterUser регистрирует юзера.
func (a *JwtAuther) RegisterUser(ctx context.Context) (types.UserID, string, error) {
	userID, err := a.usersStorage.Create(ctx)
	if err != nil {
		return users.GuestID, "", err
	}

	token, err := a.createAuthToken(userID)
	if err != nil {
		return users.GuestID, "", err
	}

	return userID, token, nil
}

// FromUserContext достаёт юзера контекста.
func (a *JwtAuther) FromUserContext(ctx context.Context) (types.UserID, error) {
	userID, ok := ctx.Value(ctxUserID).(types.UserID)
	if !ok {
		return users.GuestID, errors.New("user is not set")
	}

	return userID, nil
}

func (a *JwtAuther) createAuthToken(userID types.UserID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{},
		UserID:           int64(userID),
	})

	tokenString, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// NewAuther достаёт юзера контекста.
func NewAuther(usersStorage users.Storage, secretKey string) *JwtAuther {
	return &JwtAuther{
		usersStorage: usersStorage,
		secretKey:    secretKey,
	}
}

// NewUserContext - создаёт контекст с данными юзера.
func NewUserContext(ctx context.Context, userID types.UserID) context.Context {
	return context.WithValue(ctx, ctxUserID, userID)
}
