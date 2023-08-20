package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/microboard/database"
	"net/http"
	"strconv"
)

func GetThreads(db *database.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		boardCode := ctx.Query("boardCode")

		limit, err := strconv.Atoi(ctx.Query("limit"))
		if err != nil {
			ctx.Error(err)
			return
		}

		offset, err := strconv.Atoi(ctx.Query("offset"))
		if err != nil {
			ctx.Error(err)
			return
		}

		if limit > 100 {
			limit = 100
		}

		threads, err := db.GetThreads(boardCode, limit, offset)

		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, threads)
	}
}

type NewThread struct {
	// Text with markup.
	Body      string `json:"body" binding:"required"`
	BoardCode string `json:"boardCode" binding:"required"`
}

func CreateThread(db *database.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestData NewThread
		if err := ctx.ShouldBindJSON(&requestData); err != nil {
			ctx.Error(err)
			return
		}

		thread := database.Thread{
			BoardCode: requestData.BoardCode,
		}
		firstPost := database.Post{
			Body: requestData.Body,
		}

		if err := db.CreateThread(&thread, &firstPost, []database.File{}); err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, thread)
	}
}
