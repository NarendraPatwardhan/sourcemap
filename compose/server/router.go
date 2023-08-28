package server

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"machinelearning.one/sourcemap/compose/logger"
	"machinelearning.one/sourcemap/compose/static"
	"machinelearning.one/sourcemap/frontend"
)

// Internal function to create a production ready instance of gin router.
func createRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	return engine
}

// Creates and runs a http server suitable for SPA applications.
func Run(ctx context.Context, port uint, decoupled bool, fns ...Func) error {
	lg := logger.Get(ctx)

	// Initiate a custom instance of gin router.
	router := createRouter()
	if decoupled {
		frontendPort := os.Getenv("SOURCEMAP_FRONTEND_PORT")
		if frontendPort == "" {
			lg.Fatal().Msg("SOURCEMAP_FRONTEND_PORT not set")
		}
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{fmt.Sprintf("http://localhost:%s", frontendPort)},
			AllowMethods:     []string{"GET", "POST"},
			AllowHeaders:     []string{"*"},
			AllowCredentials: true,
		}))
	}

	// Create a group for all api routes.
	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
			})
		})
		for _, fn := range fns {
			api.Handle(fn.HTTPMethod, fn.RelativePath, fn.Handlers...)
		}
	}

	if !decoupled {
		// Load the .html files as templates to be served.
		tmpl := template.Must(
			template.New("").
				Delims("{{", "}}").
				Funcs(router.FuncMap).
				ParseFS(frontend.Content, "dist/*.html"),
		)
		router.SetHTMLTemplate(tmpl)
		// Serve the index.html file for all routes except the api ones.
		router.NoRoute(func(c *gin.Context) {
			if !strings.HasPrefix(c.Request.RequestURI, "/api") {
				c.HTML(http.StatusOK, "index.html", nil)
			}
		})

		// Load and serve static content from root.
		content, err := static.New(frontend.Content, "dist")
		if err != nil {
			return err
		}
		router.Use(static.Serve("/", *content))
	}

	// Create a new http server and attach the router to it.
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	// Start the server in a separate goroutine.
	go func() {
		srv.ListenAndServe()
	}()

	// Await for context cancellation.
	<-ctx.Done()
	srv.Shutdown(ctx)

	return nil
}
