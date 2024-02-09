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

	"github.com/ordinary-dev/microboard/config"
	"github.com/ordinary-dev/microboard/database"
	"github.com/ordinary-dev/microboard/storage"
)

func CreateThread(db *database.DB, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.Error(err)
			return
		}

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

		thread := database.Thread{
			BoardCode: boardCode[0],
		}
		firstPost := database.Post{
			Body: body[0],
		}
		dbFiles := make([]database.File, 0)

		files, ok := form.File["files"]
		if ok {
			for _, fileHeader := range files {
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

				dbFile := database.File{
					Filepath: filepath,
					Name:     fileHeader.Filename,
					Size:     uint64(fileHeader.Size),
					MimeType: mimetype,
				}
				dbFiles = append(dbFiles, dbFile)
			}
		}

		if err := db.CreateThread(&thread, &firstPost, dbFiles); err != nil {
			ctx.Error(err)
			return
		}

		go storage.GeneratePreviewsForPost(db, cfg, firstPost.ID)

		board, err := db.GetBoard(thread.BoardCode)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.Redirect(http.StatusFound, fmt.Sprintf("/boards/%v", board.Code))
	}
}

func ShowThread(db *database.DB, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		threadID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.Error(err)
			return
		}

		posts, err := db.GetPostsFromThread(threadID)
		if err != nil {
			ctx.Error(err)
			return
		}

		thread, err := db.GetThread(threadID)
		if err != nil {
			ctx.Error(err)
			return
		}

		render(ctx, cfg, http.StatusOK, "thread.html.tmpl", gin.H{
			"posts":     posts,
			"threadID":  threadID,
			"boardCode": thread.BoardCode,
		})
	}
}
