package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/ordinary-dev/microboard/config"
)

var (
	DB *pgxpool.Pool
)

func GetDatabaseConnection(cfg *config.Config) (err error) {
	logrus.Debug("Connecting to the database")

	url := getDbUrl(cfg, "postgres")

	DB, err = pgxpool.New(context.Background(), url)
	return err
}

func getDbUrl(cfg *config.Config, schema string) string {
	url := fmt.Sprintf("%v://%v", schema, cfg.DbUser)
	if cfg.DbPassword != nil {
		url += ":" + *cfg.DbPassword
	}

	url += fmt.Sprintf("@/%v?host=%v", cfg.DbName, cfg.DbHost)
	if cfg.DbPort != nil {
		url += fmt.Sprintf("&port=%v", *cfg.DbPort)
	}

	return url
}
