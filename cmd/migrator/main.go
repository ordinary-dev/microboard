package main

import (
	"github.com/sirupsen/logrus"

	"github.com/ordinary-dev/microboard/config"
	"github.com/ordinary-dev/microboard/db"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	if err = db.Migrate(cfg); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Migrations have been applied")
}
