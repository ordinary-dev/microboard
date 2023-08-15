package views

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/microboard/src/database"
	"net/http"
	"strconv"
)

var (
	ErrEmptyCode        = errors.New("board code is empty")
	ErrEmptyName        = errors.New("board name is empty")
	ErrEmptyDescription = errors.New("board description is empty")
	ErrInvalidPageLimit = errors.New("page limit is invalid")
	ErrInvalidBumpLimit = errors.New("bump limit is invalid")
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
			"pageCount": pageCount,
		})
	}
}

func CreateBoard(db *database.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var board database.Board
		if err := ctx.ShouldBind(&board); err != nil {
			ctx.Error(err)
			return
		}

		if board.Code == "" {
			ctx.Error(ErrEmptyCode)
			return
		}

		if board.Name == "" {
			ctx.Error(ErrEmptyName)
			return
		}

		if board.Description == "" {
			ctx.Error(ErrEmptyDescription)
			return
		}

		if board.BumpLimit < 0 {
			ctx.Error(ErrInvalidBumpLimit)
			return
		}

		if board.PageLimit < 0 {
			ctx.Error(ErrInvalidPageLimit)
			return
		}

		if err := db.CreateBoard(&board); err != nil {
			ctx.Error(err)
			return
		}

		ctx.Redirect(http.StatusFound, "/admin-panel")
	}
}

func UpdateBoard(db *database.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var board database.Board
		if err := ctx.ShouldBind(&board); err != nil {
			ctx.Error(err)
			return
		}

		if board.Code == "" {
			ctx.Error(ErrEmptyCode)
			return
		}

		if board.Name == "" {
			ctx.Error(ErrEmptyName)
			return
		}

		if board.Description == "" {
			ctx.Error(ErrEmptyDescription)
			return
		}

		if board.BumpLimit < 0 {
			ctx.Error(ErrInvalidBumpLimit)
			return
		}

		if board.PageLimit < 0 {
			ctx.Error(ErrInvalidPageLimit)
			return
		}

		if err := db.UpdateBoard(&board); err != nil {
			ctx.Error(err)
			return
		}

		ctx.Redirect(http.StatusFound, "/admin-panel")
	}
}
