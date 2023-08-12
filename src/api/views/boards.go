package views

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/microboard/src/database"
	"net/http"
	"strconv"
)

func ShowBoard(db *database.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		boardCode := ctx.Param("code")

		page, err := strconv.Atoi(ctx.Query("page"))
		if err != nil {
			page = 1
		}
		if page < 1 {
			page = 1
		}

		limit := 10
		offset := limit * (page - 1)

		board, err := db.GetBoard(boardCode)
		if err != nil {
			ctx.Error(err)
			return
		}

		threads, err := db.GetThreads(boardCode, limit, offset)
		if err != nil {
			ctx.Error(err)
			return
		}

		threadCount, err := db.GetThreadCount(boardCode)
		if err != nil {
			ctx.Error(err)
			return
		}

		pageLimit := int64(board.PageLimit)
		pageCount := threadCount/10 + 1

		if pageLimit == 0 {
			pageLimit = pageCount
		}

		if pageCount > pageLimit {
			pageCount = pageLimit
		}

		ctx.HTML(http.StatusOK, "board.html.tmpl", gin.H{
			"board":     board,
			"threads":   threads,
			"pageCount": pageLimit,
		})
	}
}
