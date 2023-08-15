package views

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/microboard/src/database"
	"net/http"
)

func ShowAdminPanel(db *database.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		boards, err := db.GetBoards()
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.HTML(http.StatusOK, "admin.html.tmpl", gin.H{
			"boards": boards,
		})
	}
}
