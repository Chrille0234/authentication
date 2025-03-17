package api

import (
	"github.com/chrille0234/auth/database"
	"github.com/chrille0234/auth/database/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterHandler(db *sqlx.DB) func(*gin.Context) {
  return func(ctx *gin.Context) {
    err := ctx.Request.ParseForm()
    if err != nil {
      ctx.Writer.Write([]byte("Invalid credentials"))
      ctx.AbortWithStatus(400)
      return
    }

    email := ctx.Request.FormValue("email")
    password := ctx.Request.FormValue("password")
    firstName := ctx.Request.FormValue("first_name")
    lastName := ctx.Request.FormValue("last_name")

    if email == "" || password == "" || firstName == "" || lastName == "" {
      ctx.Writer.Write([]byte("Invalid credentials"))
      ctx.AbortWithStatus(400)
      return
    }

    account := models.User{
      Email: email,
      Password: password,
      FirstName: firstName,
      LastName: lastName,
    }

    err = account.HashAndSalt()

    if err != nil {
      ctx.Writer.Write([]byte("Error creating account, try again later."))
      ctx.AbortWithStatus(500)
      return
    }

    err = database.CreateUser(account, db)

    if err != nil {
      ctx.Writer.Write([]byte("Error creating account, try again later."))
      ctx.AbortWithStatus(500)
      return
    }

    createTokenAndSetCookie(db, email, ctx)

	ctx.Header("HX-Redirect", "/profile")
	ctx.Status(200)
  }
}
