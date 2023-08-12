package views

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
	"github.com/ordinary-dev/microboard/src/config"
	"github.com/ordinary-dev/microboard/src/database"
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

				fileKind, err := filetype.Match(buf)
				if err != nil || fileKind == filetype.Unknown {
					fileKind.MIME.Value = "application/octet-stream"
					fileKind.Extension = "bin"
				}

				hasher := sha256.New()
				hasher.Write(buf)

				filepath := fmt.Sprintf("%x.%v", hasher.Sum(nil), fileKind.Extension)
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
					MimeType: fileKind.MIME.Value,
				}
				dbFiles = append(dbFiles, dbFile)
			}
		}

		err = db.CreatePost(&post, dbFiles)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.Redirect(http.StatusFound, fmt.Sprintf("/threads/%v", threadIDUint))
	}
}
