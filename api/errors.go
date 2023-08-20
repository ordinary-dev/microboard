package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func APIErrorHandler(ctx *gin.Context) {
	ctx.Next()

	for idx, err := range ctx.Errors {
		if idx == 0 {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		logrus.Debugf("%v", err)
	}
}

func HtmlErrorHandler(ctx *gin.Context) {
	ctx.Next()

	for idx, err := range ctx.Errors {
		if idx == 0 {
			ctx.HTML(http.StatusInternalServerError, "error.html.tmpl", gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
		}

		logrus.Debugf("%v", err)
	}
}