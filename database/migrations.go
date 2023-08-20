package database

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ordinary-dev/microboard/config"
	"github.com/sirupsen/logrus"
	"strings"
)

func Migrate(cfg *config.Config) error {
	logrus.Debug("Applying migrations")

	m, err := migrate.New(
		"file://migrations",
		// Replace postgres://... with pgx5://...
		strings.Replace(cfg.DatabaseUrl, "postgres", "pgx5", 1),
	)
	if err != nil {
		return err
	}

	return m.Up()
}
