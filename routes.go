package main

import (
	"github.com/a-h/templ"
	"github.com/chrille0234/auth/api"
	"github.com/chrille0234/auth/views/index"
	"github.com/chrille0234/auth/views/login"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func SetupAllRoutes(router *gin.Engine, db *sqlx.DB) {
	// UNPROTECTED ROUTES
	unprotected := router.Group("/")
    setupUnprotectedRoutes(unprotected)

	// PROTECTED ROUTES
	protected := router.Group("/")
	protected.Use(AuthMiddleware(db))
	protected.Use(RateLimiteMiddleware())
    setupProtectedRoutes(protected)

	// API ROUTES
	apiGroup := router.Group("/api")
    setupAPIRoutes(apiGroup, db)

	// COMPONENT ROUTES
	components := router.Group("/components")
    setupComponentRoutes(components)

}

// UNPROTECTED ROUTES
func setupUnprotectedRoutes(group *gin.RouterGroup) {
	group.GET("/", func(ctx *gin.Context) {
		handler := templ.Handler(index.Index())
		handler.ServeHTTP(ctx.Writer, ctx.Request)
	})

	group.GET("/login", func(ctx *gin.Context) {
		handler := templ.Handler(login.Login())
		handler.ServeHTTP(ctx.Writer, ctx.Request)
	})
}

// PROTECTED ROUTES
func setupProtectedRoutes(group *gin.RouterGroup) {
	group.GET("/profile", api.ProfileHandler)
}

// API ROUTES
func setupAPIRoutes(group *gin.RouterGroup, db *sqlx.DB) {
	group.POST("login", api.LoginHandler(db))
	group.POST("register", api.RegisterHandler(db))
}

// COMPONENT ROUTES
func setupComponentRoutes(group *gin.RouterGroup) {
	group.GET("/loginForm", func(ctx *gin.Context) {
		handler := templ.Handler(login.LoginForm())
		handler.ServeHTTP(ctx.Writer, ctx.Request)
	})

	group.GET("/registerForm", func(ctx *gin.Context) {
		handler := templ.Handler(login.RegisterForm())
		handler.ServeHTTP(ctx.Writer, ctx.Request)
	})
}


