package database

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ordinary-dev/microboard/config"
	"github.com/sirupsen/logrus"
)

func Migrate(cfg *config.Config) error {
	logrus.Debug("Applying migrations")

	m, err := migrate.New(
		"file://database/migrations",
		getDbUrl(cfg, "pgx5"),
	)
	if err != nil {
		return err
	}

	return m.Up()
}
