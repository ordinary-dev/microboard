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

		if err := db.CreateThread(&thread, &firstPost, dbFiles); err != nil {
			ctx.Error(err)
			return
		}

		board, err := db.GetBoard(thread.BoardCode)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.Redirect(http.StatusFound, fmt.Sprintf("/boards/%v", board.Code))
	}
}

func ShowThread(db *database.DB) gin.HandlerFunc {
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

		ctx.HTML(http.StatusOK, "thread.html.tmpl", gin.H{
			"posts":     posts,
			"threadID":  threadID,
			"boardCode": thread.BoardCode,
		})
	}
}