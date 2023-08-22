package views

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/microboard/config"
	"github.com/ordinary-dev/microboard/database"
	"net/http"
)

func ShowMainPage(db *database.DB, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		boards, err := db.GetBoards()
		if err != nil {
			ctx.Error(err)
			return
		}

		Render(ctx, cfg, http.StatusOK, "index.html.tmpl", gin.H{
			"boards": boards,
		})
	}
}
