package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/microboard/config"
	"github.com/ordinary-dev/microboard/database"
	"net/http"
)

func AuthorizationMiddleware(db *database.DB, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("microboard-token")
		if err != nil {
			token = ctx.GetHeader("Authorization")
		}

		if _, err := db.ValidateAccessToken(token); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
		}
	}
}
