package frontend

import (
	"html/template"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ordinary-dev/microboard/config"
)

// Set up urls for the frontend.
func ConfigureFrontend(engine *gin.Engine, cfg *config.Config) {
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

	frontend.GET("/", ShowMainPage(cfg))
	frontend.GET("/boards/:code", ShowBoard(cfg))

	frontend.POST("/threads", CreateThread(cfg))
	frontend.GET("/threads/:id", ShowThread(cfg))

	frontend.POST("/posts", CreatePost(cfg))

	frontend.GET("/login", ShowLoginForm(cfg))
	frontend.POST("/login", Authenticate(cfg))

	// Secret pages
	protectedFrontend := frontend.Group("/")
	protectedFrontend.Use(AuthorizationMiddleware(cfg))

	protectedFrontend.GET("/admin-panel", ShowAdminPanel(cfg))

	protectedFrontend.POST("/boards", CreateBoard(cfg))
	protectedFrontend.POST("/boards/:code", UpdateBoard())
}
