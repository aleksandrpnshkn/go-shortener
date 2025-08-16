package store

import (
	"context"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
)

type ShortenedURL struct {
	Code        string
	OriginalURL string
}

type Storage interface {
	Ping(ctx context.Context) error

	Set(ctx context.Context, url ShortenedURL) (storedURL ShortenedURL, hasConflict bool, err error)

	SetMany(ctx context.Context, urls map[string]ShortenedURL) (storedURLs map[string]ShortenedURL, hasConflict bool, err error)

	Get(ctx context.Context, code string) (originalURL string, isFound bool)

	Close() error
}

func NewStorage(
	ctx context.Context,
	databaseDSN string,
	fileStoragePath string,
	logger *zap.Logger,
) (Storage, error) {
	var storage Storage

	storage, err := NewSQLStorage(ctx, databaseDSN)
	if err == nil {
		err = runMigrations(databaseDSN)
		if err != nil {
			return nil, errors.New("failed to run migrations")
		}

		return storage, nil
	}

	logger.Warn("failed to init SQL storage", zap.Error(err))

	storage, err = NewFileStorage(fileStoragePath)
	if err != nil {
		return nil, errors.New("failed to init file storage")
	}

	return storage, nil
}

//go:embed migrations/*.sql
var migrationsFiles embed.FS

func runMigrations(databaseDSN string) error {
	sourceDriver, err := iofs.New(migrationsFiles, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, databaseDSN)
	if err != nil {
		return err
	}

	err = m.Up()
	if !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
