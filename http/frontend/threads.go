package frontend

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ordinary-dev/microboard/config"
	"github.com/ordinary-dev/microboard/db/boards"
	"github.com/ordinary-dev/microboard/db/captchas"
	"github.com/ordinary-dev/microboard/db/files"
	"github.com/ordinary-dev/microboard/db/posts"
	"github.com/ordinary-dev/microboard/db/threads"
	"github.com/ordinary-dev/microboard/storage"
)

func CreateThread(cfg *config.Config) gin.HandlerFunc {
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

		isCaptchaValid, err := captchas.ValidateCaptcha(ctx, captchaID, captchaAnswer[0])
		if err != nil {
			ctx.Error(err)
			return
		}

		if !isCaptchaValid {
			ctx.Error(errors.New("captcha is invalid"))
			return
		}

		// Create thread
		boardCode, ok := form.Value["boardCode"]
		if !ok || len(boardCode) < 1 {
			ctx.Error(errors.New("boardCode is undefined"))
			return
		}

		body, ok := form.Value["body"]
		if !ok || len(body) < 1 {
			ctx.Error(errors.New("body is undefined"))
			return
		}

		thread := threads.Thread{
			BoardCode: boardCode[0],
		}
		firstPost := posts.Post{
			Body: body[0],
		}
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

				filepath, mimetype, err := storage.SaveBuffer(cfg, buf, boardCode[0])
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

		if err := threads.CreateThread(&thread, &firstPost, dbFiles); err != nil {
			ctx.Error(err)
			return
		}

		go storage.GeneratePreviewsForPost(cfg, firstPost.ID)

		board, err := boards.GetBoard(thread.BoardCode)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.Redirect(http.StatusFound, fmt.Sprintf("/boards/%v", board.Code))
	}
}

func ShowThread(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		threadID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.Error(err)
			return
		}

		posts, err := posts.GetPostsFromThread(threadID)
		if err != nil {
			ctx.Error(err)
			return
		}

		thread, err := threads.GetThread(threadID)
		if err != nil {
			ctx.Error(err)
			return
		}

		captcha, err := captchas.CreateCaptcha(ctx)
		if err != nil {
			ctx.Error(err)
			return
		}

		render(ctx, cfg, http.StatusOK, "thread.html.tmpl", gin.H{
			"posts":     posts,
			"threadID":  threadID,
			"boardCode": thread.BoardCode,
			"captchaID": captcha.ID.String(),
		})
	}
}
