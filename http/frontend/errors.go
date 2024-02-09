package frontend

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ordinary-dev/microboard/config"
)

func HtmlErrorHandler(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		for idx, err := range ctx.Errors {
			if idx == 0 {
				render(ctx, cfg, http.StatusInternalServerError, "error.html.tmpl", gin.H{
					"error": err.Error(),
				})
				ctx.Abort()
			}

			logrus.Debugf("%v", err)
		}
	}
}
