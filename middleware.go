package main

import (
	"time"

	"github.com/JGLTechnologies/gin-rate-limit"
	"github.com/chrille0234/auth/database"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
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

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

func RateLimiteMiddleware() gin.HandlerFunc {
	// limit to 5 requests per second
	store := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: redis.NewClient(&redis.Options{
			Addr: "localhost:7680",
		}),
		Rate:  time.Second,
		Limit: 5,
	})
	mw := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})

	return mw
}
