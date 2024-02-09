package frontend

import (
	"github.com/gin-gonic/gin"

	"github.com/ordinary-dev/microboard/config"
)

// Fill in required parameters and call gin.Context.HTML().
func render(ctx *gin.Context, cfg *config.Config, status int, template string, params gin.H) {
	params["plausibleSrc"] = cfg.PlausibleAnalyticsSrc
	params["domain"] = cfg.Domain

	ctx.HTML(status, template, params)
}

// Show error page.
// Calls gin.Context.Abort().
func renderError(ctx *gin.Context, cfg *config.Config, status int, err error) {
	render(ctx, cfg, status, "error.html.tmpl", gin.H{
		"error": err.Error(),
	})
	ctx.Abort()
}
