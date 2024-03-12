package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ordinary-dev/microboard/config"
	"github.com/ordinary-dev/microboard/db/boards"
	"github.com/ordinary-dev/microboard/storage"
)

func CreateBoard(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var newBoard boards.Board
		if err := ctx.ShouldBindJSON(&newBoard); err != nil {
			ctx.Error(err)
			return
		}

		if err := boards.CreateBoard(&newBoard); err != nil {
			ctx.Error(err)
			return
		}

		if err := storage.CreateDirs(cfg); err != nil {
			logrus.Fatal(err)
		}

		ctx.JSON(http.StatusCreated, newBoard)
	}
}

func GetBoards() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		boardList, err := boards.GetBoards()
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, boardList)
	}
}

func GetBoard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		board, err := boards.GetBoard(ctx.Param("code"))
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, board)
	}
}

func UpdateBoard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var updatedBoard boards.Board
		if err := ctx.ShouldBindJSON(&updatedBoard); err != nil {
			ctx.Error(err)
			return
		}

		if err := boards.UpdateBoard(&updatedBoard); err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, updatedBoard)
	}
}

func DeleteBoard(ctx *gin.Context) {
	code := ctx.Param("code")

	if err := boards.DeleteBoard(ctx, code); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
