package views

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/microboard/config"
	"github.com/ordinary-dev/microboard/database"
	"net/http"
)

func ShowLoginForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "auth.html.tmpl", gin.H{})
}

type AuthForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func Authenticate(db *database.DB, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var authForm AuthForm
		if err := ctx.ShouldBind(&authForm); err != nil {
			ctx.Error(err)
			return
		}

		admin, err := db.VerifyPassword(authForm.Username, authForm.Password)
		if err != nil {
			ctx.Error(err)
			return
		}

		token, err := db.GetAccessToken(admin.ID)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.SetCookie("microboard-token", token.Value, 60*60*24*7, "/", "", cfg.IsProduction, true)
		ctx.Redirect(http.StatusFound, "/admin-panel")
	}
}

func AuthenticationMiddleware(db *database.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("microboard-token")
		if err != nil {
			ctx.HTML(http.StatusUnauthorized, "error.html.tmpl", gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
			return
		}

		if _, err := db.ValidateAccessToken(token); err != nil {
			ctx.HTML(http.StatusUnauthorized, "error.html.tmpl", gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
		}
	}
}
