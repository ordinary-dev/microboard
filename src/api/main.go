package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/microboard/src/api/views"
	"github.com/ordinary-dev/microboard/src/config"
	"github.com/ordinary-dev/microboard/src/database"
	"html/template"
)

func GetAPIEngine(db *database.DB, cfg *config.Config) *gin.Engine {
	if cfg.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	engine.SetFuncMap(template.FuncMap{
		"loop": func(from, to int64) <-chan int64 {
			ch := make(chan int64)
			go func() {
				for i := from; i <= to; i++ {
					ch <- i
				}
				close(ch)
			}()
			return ch
		},
	})
	engine.LoadHTMLGlob("templates/*")
	engine.Static("/assets", "./assets")
	engine.Static("/uploads", cfg.UploadDir)

	// Public pages
	frontend := engine.Group("/")
	frontend.Use(HtmlErrorHandler)
	frontend.GET("/", views.ShowMainPage(db))
	frontend.GET("/boards/:code", views.ShowBoard(db))
	frontend.POST("/threads", views.CreateThread(db, cfg))
	frontend.GET("/threads/:id", views.ShowThread(db))
	frontend.POST("/posts", views.CreatePost(db, cfg))
	frontend.GET("/login", views.ShowLoginForm)
	frontend.POST("/login", views.Authenticate(db, cfg))

	// Secret pages
	protectedFrontend := frontend.Group("/")
	protectedFrontend.Use(views.AuthenticationMiddleware(db))
	protectedFrontend.GET("/admin-panel", views.ShowAdminPanel(db))
	protectedFrontend.POST("/boards", views.CreateBoard(db))
	protectedFrontend.POST("/boards/:code", views.UpdateBoard(db))

	// JSON API
	v0 := engine.Group("/api/v0")
	v0.Use(APIErrorHandler)

	v0.POST("/users", CreateUser(db))
	v0.POST("/users/token", GetAccessToken(db))

	// Boards
	v0.POST("/boards", CreateBoard(db))
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

	return engine
}
