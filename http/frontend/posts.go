package frontend

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ordinary-dev/microboard/config"
	dbcaptchas "github.com/ordinary-dev/microboard/db/captchas"
	"github.com/ordinary-dev/microboard/db/files"
	"github.com/ordinary-dev/microboard/db/posts"
	"github.com/ordinary-dev/microboard/db/threads"
	"github.com/ordinary-dev/microboard/storage"
)

func CreatePost(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.Error(err)
			return
		}

		// Validate captcha
		captchaIDText, ok := form.Value["captchaID"]
		if !ok || len(captchaIDText) < 1 {
			ctx.Error(errors.New("captcha id is undefined"))
			return
		}

		captchaID, err := uuid.Parse(captchaIDText[0])
		if err != nil {
			ctx.Error(err)
			return
		}

		captchaAnswer, ok := form.Value["answer"]
		if !ok || len(captchaAnswer) < 1 {
			ctx.Error(errors.New("captcha answer is undefined"))
			return
		}

		isCaptchaValid, err := dbcaptchas.ValidateCaptcha(ctx, captchaID, captchaAnswer[0])
		if err != nil {
			ctx.Error(err)
			return
		}

		if !isCaptchaValid {
			ctx.Error(errors.New("captcha is invalid"))
			return
		}

		// Get thread info
		threadID, ok := form.Value["threadID"]
		if !ok || len(threadID) < 1 {
			ctx.Error(errors.New("threadID is undefined"))
			return
		}

		threadIDUint, err := strconv.ParseUint(threadID[0], 10, 64)
		if err != nil {
			ctx.Error(err)
			return
		}

		thread, err := threads.GetThread(threadIDUint)
		if err != nil {
			ctx.Error(err)
			return
		}

		// Get text
		body, ok := form.Value["body"]
		if !ok || len(body) < 1 {
			ctx.Error(errors.New("body is undefined"))
			return
		}

		post := posts.Post{
			ThreadID:  threadIDUint,
			Body:      body[0],
			CreatedAt: time.Now(),
		}

		// Save files
		dbFiles := make([]files.File, 0)

		fileList, ok := form.File["files"]
		if ok {
			for _, fileHeader := range fileList {
				file, err := fileHeader.Open()
				if err != nil {
					ctx.Error(err)
					return
				}

				buf, err := io.ReadAll(file)
				if err != nil {
					ctx.Error(err)
					return
				}

				filepath, mimetype, err := storage.SaveBuffer(cfg, buf, thread.BoardCode)
				if err != nil {
					ctx.Error(err)
					return
				}

				fullFilePath := path.Join(cfg.UploadDir, filepath)
				if _, err := os.Stat(fullFilePath); errors.Is(err, os.ErrNotExist) {
					if err = os.WriteFile(fullFilePath, buf, 0644); err != nil {
						ctx.Error(err)
						return
					}
				}

				dbFile := files.File{
					Filepath: filepath,
					Name:     fileHeader.Filename,
					Size:     uint64(fileHeader.Size),
					MimeType: mimetype,
				}
				dbFiles = append(dbFiles, dbFile)
			}
		}

		// Save post
		err = posts.CreatePost(&post, dbFiles)
		if err != nil {
			ctx.Error(err)
			return
		}

		go storage.GeneratePreviewsForPost(cfg, post.ID)

		ctx.Redirect(http.StatusFound, fmt.Sprintf("/threads/%v", threadIDUint))
	}
}
