package frontend

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ordinary-dev/microboard/config"
	"github.com/ordinary-dev/microboard/db/boards"
)

func ShowAdminPanel(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		boardList, err := boards.GetBoards()
		if err != nil {
			ctx.Error(err)
			return
		}

		render(ctx, cfg, http.StatusOK, "admin.html.tmpl", gin.H{
			"boards": boardList,
		})
	}
}
