package main

import (
	"fmt"
	"github.com/davidbyttow/govips/v2/vips"
	"github.com/ordinary-dev/microboard/src/api"
	"github.com/ordinary-dev/microboard/src/config"
	"github.com/ordinary-dev/microboard/src/database"
	"github.com/ordinary-dev/microboard/src/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	vips.Startup(nil)
	defer vips.Shutdown()

	cfg, err := config.GetConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	logLevel := cfg.GetLogLevel()
	logrus.SetLevel(logLevel)
	logrus.Infof("Setting log level to %v", logLevel)

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

	if err := storage.CreateDirs(cfg, db); err != nil {
		logrus.Fatal(err)
	}

	go storage.GenerateMissingPreviews(db, cfg)

	engine := api.GetAPIEngine(db, cfg)
	engine.Run(fmt.Sprintf("%v:%v", cfg.Addr, cfg.Port))
}
