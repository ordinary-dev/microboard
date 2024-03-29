package frontend

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ordinary-dev/microboard/config"
	"github.com/ordinary-dev/microboard/db/boards"
	"github.com/ordinary-dev/microboard/db/captchas"
	"github.com/ordinary-dev/microboard/db/threads"
	"github.com/ordinary-dev/microboard/storage"
)

var (
	ErrEmptyCode        = errors.New("board code is empty")
	ErrEmptyName        = errors.New("board name is empty")
	ErrEmptyDescription = errors.New("board description is empty")
	ErrInvalidPageLimit = errors.New("page limit is invalid")
	ErrInvalidBumpLimit = errors.New("bump limit is invalid")
)

func ShowBoard(cfg *config.Config) gin.HandlerFunc {
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

		board, err := boards.GetBoard(boardCode)
		if err != nil {
			ctx.Error(err)
			return
		}

		threadList, err := threads.GetThreads(boardCode, limit, offset)
		if err != nil {
			ctx.Error(err)
			return
		}

		threadCount, err := threads.GetThreadCount(boardCode)
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

		captcha, err := captchas.CreateCaptcha(ctx)
		if err != nil {
			ctx.Error(err)
			return
		}

		render(ctx, cfg, http.StatusOK, "board.html.tmpl", gin.H{
			"board":     board,
			"threads":   threadList,
			"pageCount": pageCount,
			"captchaID": captcha.ID.String(),
		})
	}
}

func CreateBoard(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var board boards.Board
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

		if err := boards.CreateBoard(&board); err != nil {
			ctx.Error(err)
			return
		}

		if err := storage.CreateDirs(cfg); err != nil {
			logrus.Fatal(err)
		}

		ctx.Redirect(http.StatusFound, "/admin-panel")
	}
}

func UpdateBoard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var board boards.Board
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

		if err := boards.UpdateBoard(&board); err != nil {
			ctx.Error(err)
			return
		}

		ctx.Redirect(http.StatusFound, "/admin-panel")
	}
}
