package main

import (
	"fmt"
	"github.com/ordinary-dev/microboard/src/api"
	"github.com/ordinary-dev/microboard/src/config"
	"github.com/ordinary-dev/microboard/src/database"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	logLevel := cfg.GetLogLevel()
	logrus.SetLevel(logLevel)
	logrus.Infof("Setting log level to %v", logLevel)

	if err := os.MkdirAll(cfg.UploadDir, os.ModePerm); err != nil {
		logrus.Fatal(err)
	}

	if err = database.Migrate(cfg); err != nil {
		logrus.Errorf("Migrate: %v", err)
	}

	db, err := database.GetDatabaseConnection(cfg)
	if err != nil {
		logrus.Fatalln(err)
	}

	if err := db.CreateDefaultUser(cfg); err != nil {
		logrus.Error(err)
	}

	engine := api.GetAPIEngine(db, cfg)
	engine.Run(fmt.Sprintf("%v:%v", cfg.Addr, cfg.Port))
}
