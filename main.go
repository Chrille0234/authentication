package main

import (
	"github.com/a-h/templ"
	"github.com/chrille0234/auth/api"
	"github.com/chrille0234/auth/database"
	"github.com/chrille0234/auth/views/index"
	"github.com/chrille0234/auth/views/login"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db := database.ConnectAndSeed("db.sqlite3")
	defer db.Close()

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.LoadHTMLGlob("./static/*.css")
	router.Static("/static", "./static")

	// UNPROTECTED ROUTES
	unprotected := router.Group("/")
	{
		unprotected.GET("/", func(ctx *gin.Context) {
			handler := templ.Handler(index.Index())
			handler.ServeHTTP(ctx.Writer, ctx.Request)
		})

		unprotected.GET("/login", func(ctx *gin.Context) {
			handler := templ.Handler(login.Login())
			handler.ServeHTTP(ctx.Writer, ctx.Request)
		})
	}

	// PROTECTED ROUTES
	protected := router.Group("/")
	protected.Use(AuthMiddleware(db))
	{
		protected.GET("/profile", api.ProfileHandler)
	}

	// API ROUTES
	apiGroup := router.Group("/api")
	{
		apiGroup.POST("login", api.LoginHandler(db))
        apiGroup.POST("register", api.RegisterHandler(db))
	}

	// COMPONENT ROUTES
	components := router.Group("/components")
	{
		components.GET("/loginForm", func(ctx *gin.Context) {
			handler := templ.Handler(login.LoginForm())
			handler.ServeHTTP(ctx.Writer, ctx.Request)
		})

		components.GET("/registerForm", func(ctx *gin.Context) {
			handler := templ.Handler(login.RegisterForm())
			handler.ServeHTTP(ctx.Writer, ctx.Request)
		})
	}
	router.Run(":4000")
}
