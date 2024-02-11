package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	dbcaptchas "github.com/ordinary-dev/microboard/db/captchas"
	"github.com/ordinary-dev/microboard/db/files"
	"github.com/ordinary-dev/microboard/db/posts"
	"github.com/ordinary-dev/microboard/db/threads"
)

func GetThreads() gin.HandlerFunc {
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

		threadList, err := threads.GetThreads(boardCode, limit, offset)

		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, threadList)
	}
}

type NewThread struct {
	// Text with markup.
	Body          string `json:"body" binding:"required"`
	BoardCode     string `json:"boardCode" binding:"required"`
	CaptchaID     string `json:"captchaID" binding:"required"`
	CaptchaAnswer string `json:"captchaAnswer" binding:"required"`
}

func CreateThread() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestData NewThread
		if err := ctx.ShouldBindJSON(&requestData); err != nil {
			ctx.Error(err)
			return
		}

		// Validate captcha
		captchaID, err := uuid.Parse(requestData.CaptchaID)
		if err != nil {
			ctx.Error(err)
			return
		}

		isCaptchaValid, err := dbcaptchas.ValidateCaptcha(ctx, captchaID, requestData.CaptchaAnswer)
		if err != nil {
			ctx.Error(err)
			return
		}

		if !isCaptchaValid {
			ctx.Error(errors.New("captcha is invalid"))
			return
		}

		// Create thread
		thread := threads.Thread{
			BoardCode: requestData.BoardCode,
		}
		firstPost := posts.Post{
			Body: requestData.Body,
		}

		if err := threads.CreateThread(&thread, &firstPost, []files.File{}); err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, thread)
	}
}
