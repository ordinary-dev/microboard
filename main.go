package main

import (
	"fmt"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/sirupsen/logrus"

	"github.com/ordinary-dev/microboard/config"
	"github.com/ordinary-dev/microboard/db"
	"github.com/ordinary-dev/microboard/db/users"
	"github.com/ordinary-dev/microboard/http"
	"github.com/ordinary-dev/microboard/storage"
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

	if err = db.Migrate(cfg); err != nil {
		logrus.Errorf("Migrate: %v", err)
	}

	if err := db.GetDatabaseConnection(cfg); err != nil {
		logrus.Fatal(err)
	}

	if err := users.CreateDefaultUser(cfg); err != nil {
		logrus.Error(err)
	}

	if err := storage.CreateDirs(cfg); err != nil {
		logrus.Fatal(err)
	}

	go storage.GenerateMissingPreviews(cfg)

	engine := http.GetEngine(cfg)
	engine.Run(fmt.Sprintf("%v:%v", cfg.Addr, cfg.Port))
}
