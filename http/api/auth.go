package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ordinary-dev/microboard/config"
	"github.com/ordinary-dev/microboard/db/users"
)

func AuthorizationMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("microboard-token")
		if err != nil {
			token = ctx.GetHeader("Authorization")
		}

		if _, err := users.ValidateAccessToken(token); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
		}
	}
}
