package store

import (
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

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
