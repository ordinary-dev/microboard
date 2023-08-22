package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ordinary-dev/microboard/config"
	"github.com/sirupsen/logrus"
)

type DB struct {
	pool *pgxpool.Pool
}

func GetDatabaseConnection(cfg *config.Config) (*DB, error) {
	logrus.Debug("Connecting to the database")

	url := getDbUrl(cfg, "postgres")

	dbpool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}

	db := DB{
		pool: dbpool,
	}

	return &db, nil
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
