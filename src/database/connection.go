package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ordinary-dev/microboard/src/config"
	"github.com/sirupsen/logrus"
)

type DB struct {
	pool *pgxpool.Pool
}

func GetDatabaseConnection(cfg *config.Config) (*DB, error) {
	logrus.Debug("Connecting to the database")
	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseUrl)
	if err != nil {
		return nil, err
	}

	db := DB{
		pool: dbpool,
	}

	return &db, nil
}
