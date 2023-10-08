package frontend

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/microboard/config"
	"github.com/sirupsen/logrus"
	"net/http"
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
