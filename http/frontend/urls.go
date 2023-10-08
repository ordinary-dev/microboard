package frontend

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/microboard/config"
	"github.com/ordinary-dev/microboard/database"
	"html/template"
	"strings"
)

// Set up urls for the frontend.
func ConfigureFrontend(engine *gin.Engine, db *database.DB, cfg *config.Config) {
	frontend := engine.Group("/")
	frontend.Use(HtmlErrorHandler(cfg))

	// Additional functions that are used when rendering html pages.
	engine.SetFuncMap(template.FuncMap{
		// Generate a sequence of numbers.
		// Used to generate multiple elements in a loop.
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
		"hasPrefix": strings.HasPrefix,
		"isNotEmpty": func(t string) bool {
			return t != ""
		},
	})

	engine.LoadHTMLGlob("templates/*")

	frontend.GET("/", ShowMainPage(db, cfg))
	frontend.GET("/boards/:code", ShowBoard(db, cfg))

	frontend.POST("/threads", CreateThread(db, cfg))
	frontend.GET("/threads/:id", ShowThread(db, cfg))

	frontend.POST("/posts", CreatePost(db, cfg))

	frontend.GET("/login", ShowLoginForm(cfg))
	frontend.POST("/login", Authenticate(db, cfg))

	// Secret pages
	protectedFrontend := frontend.Group("/")
	protectedFrontend.Use(AuthenticationMiddleware(db, cfg))

	protectedFrontend.GET("/admin-panel", ShowAdminPanel(db, cfg))

	protectedFrontend.POST("/boards", CreateBoard(db, cfg))
	protectedFrontend.POST("/boards/:code", UpdateBoard(db))
}
