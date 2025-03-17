package main

import (
	"github.com/chrille0234/auth/database"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

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
