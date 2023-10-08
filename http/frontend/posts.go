package frontend

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/microboard/config"
	"github.com/ordinary-dev/microboard/database"
	"github.com/ordinary-dev/microboard/storage"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

func CreatePost(db *database.DB, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.Error(err)
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

		thread, err := db.GetThread(threadIDUint)
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

		post := database.Post{
			ThreadID:  threadIDUint,
			Body:      body[0],
			CreatedAt: time.Now(),
		}

		// Save files
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

				dbFile := database.File{
					Filepath: filepath,
					Name:     fileHeader.Filename,
					Size:     uint64(fileHeader.Size),
					MimeType: mimetype,
				}
				dbFiles = append(dbFiles, dbFile)
			}
		}

		// Save post
		err = db.CreatePost(&post, dbFiles)
		if err != nil {
			ctx.Error(err)
			return
		}

		go storage.GeneratePreviewsForPost(db, cfg, post.ID)

		ctx.Redirect(http.StatusFound, fmt.Sprintf("/threads/%v", threadIDUint))
	}
}
