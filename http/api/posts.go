package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	dbcaptchas "github.com/ordinary-dev/microboard/db/captchas"
	"github.com/ordinary-dev/microboard/db/files"
	"github.com/ordinary-dev/microboard/db/posts"
)

func GetPosts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		threadID, err := strconv.ParseUint(ctx.Query("threadID"), 10, 64)
		if err != nil {
			ctx.Error(err)
			return
		}

		postList, err := posts.GetPostsFromThread(threadID)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, postList)
	}
}

type NewPost struct {
	// Text with markup.
	Body          string `gorm:"not null" json:"text" binding:"required"`
	ThreadID      uint64 `json:"threadID" binding:"required"`
	CaptchaID     string `json:"captchaID" binding:"required"`
	CaptchaAnswer string `json:"captchaAnswer" binding:"required"`
}

func CreatePost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestData NewPost
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

		// Create post
		post := posts.Post{
			ThreadID:  requestData.ThreadID,
			Body:      requestData.Body,
			CreatedAt: time.Now(),
		}

		err = posts.CreatePost(&post, []files.File{})
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, post)
	}
}
