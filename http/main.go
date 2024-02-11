package http

import (
	"github.com/gin-gonic/gin"

	"github.com/ordinary-dev/microboard/config"
	"github.com/ordinary-dev/microboard/http/api"
	"github.com/ordinary-dev/microboard/http/frontend"
)

func GetEngine(cfg *config.Config) *gin.Engine {
	if cfg.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	engine.Static("/assets", "./assets")
	engine.Static("/uploads", cfg.UploadDir)
	engine.Static("/previews", cfg.PreviewDir)

	// HTML pages
	frontend.ConfigureFrontend(engine, cfg)

	// JSON API
	api.ConfigureAPI(engine, cfg)

	return engine
}
