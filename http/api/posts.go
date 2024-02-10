package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ordinary-dev/microboard/database"
	dbcaptchas "github.com/ordinary-dev/microboard/database/captchas"
)

func GetPosts(db *database.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		threadID, err := strconv.ParseUint(ctx.Query("threadID"), 10, 64)
		if err != nil {
			ctx.Error(err)
			return
		}

		posts, err := db.GetPostsFromThread(threadID)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, posts)
	}
}

type NewPost struct {
	// Text with markup.
	Body          string `gorm:"not null" json:"text" binding:"required"`
	ThreadID      uint64 `json:"threadID" binding:"required"`
	CaptchaID     string `json:"captchaID" binding:"required"`
	CaptchaAnswer string `json:"captchaAnswer" binding:"required"`
}

func CreatePost(db *database.DB) gin.HandlerFunc {
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

		isCaptchaValid, err := dbcaptchas.ValidateCaptcha(ctx, db.Pool, captchaID, requestData.CaptchaAnswer)
		if err != nil {
			ctx.Error(err)
			return
		}

		if !isCaptchaValid {
			ctx.Error(errors.New("captcha is invalid"))
			return
		}

		// Create post
		post := database.Post{
			ThreadID:  requestData.ThreadID,
			Body:      requestData.Body,
			CreatedAt: time.Now(),
		}

		err = db.CreatePost(&post, []database.File{})
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, post)
	}
}
