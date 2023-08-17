package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	// General settings.

	// Long random string used for authorization.
	SecretKey string `required:"true"`
	// Log level: debug, info, warning, error, fatal.
	LogLevel string `default:"warning"`

	// Database settings.

	// Postgres db url: "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Utc"
	DatabaseUrl string `required:"true"`

	// HTTP server settings.

	// The address on which the server will listen.
	Addr string `default:"0.0.0.0"`
	// The port on which the server will listen.
	Port int `default:"8000"`
	// Gin mode
	IsProduction bool `default:"false"`

	UploadDir  string `default:"uploads"`
	PreviewDir string `default:"previews"`

	DefaultUsername string
	DefaultPassword string
}

func GetConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logrus.Warning(err)
	}

	var config Config
	err = envconfig.Process("mb", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (cfg *Config) GetLogLevel() logrus.Level {
	switch cfg.LogLevel {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	default:
		logrus.Errorf("MB_LOGLEVEL=%v is not supported", cfg.LogLevel)
		return logrus.WarnLevel
	}
}
