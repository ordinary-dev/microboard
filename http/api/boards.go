package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ordinary-dev/microboard/config"
	"github.com/ordinary-dev/microboard/database"
	"github.com/ordinary-dev/microboard/storage"
)

func CreateBoard(db *database.DB, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var newBoard database.Board
		if err := ctx.ShouldBindJSON(&newBoard); err != nil {
			ctx.Error(err)
			return
		}

		if err := db.CreateBoard(&newBoard); err != nil {
			ctx.Error(err)
			return
		}

		if err := storage.CreateDirs(cfg, db); err != nil {
			logrus.Fatal(err)
		}

		ctx.JSON(http.StatusCreated, newBoard)
	}
}

func GetBoards(db *database.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		boards, err := db.GetBoards()
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, boards)
	}
}

func GetBoard(db *database.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		board, err := db.GetBoard(ctx.Param("code"))
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, board)
	}
}

func UpdateBoard(db *database.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var updatedBoard database.Board
		if err := ctx.ShouldBindJSON(&updatedBoard); err != nil {
			ctx.Error(err)
			return
		}

		if err := db.UpdateBoard(&updatedBoard); err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, updatedBoard)
	}
}
