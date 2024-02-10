package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ordinary-dev/microboard/database"
	dbcaptchas "github.com/ordinary-dev/microboard/database/captchas"
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
	Body          string `json:"body" binding:"required"`
	BoardCode     string `json:"boardCode" binding:"required"`
	CaptchaID     string `json:"captchaID" binding:"required"`
	CaptchaAnswer string `json:"captchaAnswer" binding:"required"`
}

func CreateThread(db *database.DB) gin.HandlerFunc {
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

		isCaptchaValid, err := dbcaptchas.ValidateCaptcha(ctx, db.Pool, captchaID, requestData.CaptchaAnswer)
		if err != nil {
			ctx.Error(err)
			return
		}

		if !isCaptchaValid {
			ctx.Error(errors.New("captcha is invalid"))
			return
		}

		// Create thread
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
