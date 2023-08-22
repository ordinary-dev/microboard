package views

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/microboard/config"
)

// Fill in required parameters and call gin.Context.HTML().
func Render(ctx *gin.Context, cfg *config.Config, status int, template string, params gin.H) {
	params["plausibleSrc"] = cfg.PlausibleAnalyticsSrc
	params["domain"] = cfg.Domain

	ctx.HTML(status, template, params)
}
