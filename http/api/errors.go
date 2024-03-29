package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
