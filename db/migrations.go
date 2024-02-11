package db

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"

	"github.com/ordinary-dev/microboard/config"
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
