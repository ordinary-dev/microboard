package frontend

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/microboard/config"
	"github.com/ordinary-dev/microboard/database"
	"net/http"
)

func ShowLoginForm(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		render(ctx, cfg, http.StatusOK, "auth.html.tmpl", gin.H{})
	}
}

type AuthForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// Path: "/login"
// Method: "POST"
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

func AuthorizationMiddleware(db *database.DB, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("microboard-token")
		if err != nil {
			renderError(ctx, cfg, http.StatusUnauthorized, err)
			return
		}

		if _, err := db.ValidateAccessToken(token); err != nil {
			renderError(ctx, cfg, http.StatusUnauthorized, err)
		}
	}
}
