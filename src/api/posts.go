package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/microboard/src/database"
	"net/http"
	"strconv"
	"time"
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
	Body     string `gorm:"not null" json:"text" binding:"required"`
	ThreadID uint64 `json:"threadID" binding:"required"`
}

func CreatePost(db *database.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestData NewPost
		if err := ctx.ShouldBindJSON(&requestData); err != nil {
			ctx.Error(err)
			return
		}

		post := database.Post{
			ThreadID:  requestData.ThreadID,
			Body:      requestData.Body,
			CreatedAt: time.Now(),
		}

		err := db.CreatePost(&post, []database.File{})
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, post)
	}
}
