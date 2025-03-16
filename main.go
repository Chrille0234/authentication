package main

import (
	"github.com/a-h/templ"
	"github.com/chrille0234/auth/api"
	"github.com/chrille0234/auth/database"
	"github.com/chrille0234/auth/views/index"
	"github.com/chrille0234/auth/views/login"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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
		apiGroup.POST("login", func(ctx *gin.Context) {
			api.LoginHandler(ctx, db)
		})
	}
	router.Run(":3000")
}

func AuthMiddleware(db *sqlx.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("auth_token")

		if err != nil {
			ctx.Writer.Write([]byte("<script>window.location.href = '/login'</script>"))
			ctx.AbortWithStatus(401)
		}

		isValid, err := database.AuthenticateToken(cookie)

		if err != nil || !isValid {
			ctx.Writer.Write([]byte("<script>window.location.href = '/login'</script>"))
			ctx.AbortWithStatus(401)
		}

		user, err := database.GetUserFromToken(cookie, db)

		if err != nil {
			ctx.Writer.Write([]byte("<script>window.location.href = '/login'</script>"))
			ctx.AbortWithStatus(401)
		}

		ctx.Set("user", user)

		ctx.Next()
	}
}
