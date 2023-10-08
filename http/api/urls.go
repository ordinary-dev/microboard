package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/microboard/config"
	"github.com/ordinary-dev/microboard/database"
)

func ConfigureAPI(engine *gin.Engine, db *database.DB, cfg *config.Config) {
	// JSON API
	v0 := engine.Group("/api/v0")
	v0.Use(APIErrorHandler)

	v0.POST("/users", CreateUser(db))
	v0.POST("/users/token", GetAccessToken(db))

	// Boards
	v0.POST("/boards", CreateBoard(db, cfg))
	v0.GET("/boards", GetBoards(db))
	v0.GET("/boards/:code", GetBoard(db))
	v0.PUT("/boards/:code", UpdateBoard(db))

	// Threads

	// Example: GET /threads?boardCode=b&limit=10&offset=0
	v0.GET("/threads", GetThreads(db))
	v0.POST("/threads", CreateThread(db))

	// Post
	v0.GET("/posts", GetPosts(db))
	v0.POST("/posts", CreatePost(db))
}