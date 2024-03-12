package api

import (
	"github.com/gin-gonic/gin"

	"github.com/ordinary-dev/microboard/config"
)

func ConfigureAPI(engine *gin.Engine, cfg *config.Config) {
	// JSON API
	v0 := engine.Group("/api/v0")
	v0.Use(APIErrorHandler)

	v0.POST("/users", CreateUser())
	v0.POST("/users/token", GetAccessToken())

	// Boards
	v0.GET("/boards", GetBoards())
	v0.GET("/boards/:code", GetBoard())

	// Threads

	// Example: GET /threads?boardCode=b&limit=10&offset=0
	v0.GET("/threads", GetThreads())
	v0.POST("/threads", CreateThread())

	// Post
	v0.GET("/posts", GetPosts())
	v0.POST("/posts", CreatePost())

	v0.GET("/captcha/:id", ShowCaptcha())

	// Protected API routes (admin-only)
	protectedAPI := v0.Group("/")
	protectedAPI.Use(AuthorizationMiddleware(cfg))
	protectedAPI.POST("/boards", CreateBoard(cfg))
	protectedAPI.PUT("/boards/:code", UpdateBoard())
	protectedAPI.DELETE("/boards/:code", DeleteBoard)
}
