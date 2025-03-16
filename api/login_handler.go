package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chrille0234/auth/database"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func LoginHandler(ctx *gin.Context, db *sqlx.DB) {
	err := ctx.Request.ParseForm()

	if err != nil {
		log.Printf("Error parsing form: %v", err)
		ctx.Writer.Write([]byte("Error processing login request"))
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	email := ctx.Request.FormValue("email")
	password := ctx.Request.FormValue("password")

	isValidPassword, err := database.AuthenticateUserPassword(email, password, db)

	if err != nil {
		log.Printf("Error authenticating user: %v", err)
		ctx.Writer.Write([]byte("Error processing login request"))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if !isValidPassword {
		ctx.Writer.Write([]byte("Invalid email or password"))
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	expires := time.Now().Add(8 * time.Hour).Unix()
	// expires is an int64, so we need to convert it to an int for the cookie
	maxAge := int(time.Until(time.Unix(expires, 0)).Seconds())
	token, err := database.CreateToken(email, expires)

	if err != nil {
		fmt.Println(err)
		ctx.Writer.Write([]byte("Error logging in, try again later"))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, err := database.GetUserUsingEmail(email, db)

	database.SetTokenInDB(token, expires, user, db)

	ctx.SetCookie("auth_token", token, maxAge, "/", "localhost", true, true)

	ctx.Header("HX-Redirect", "/profile")
	ctx.Status(200)
}
