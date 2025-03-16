package api

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/chrille0234/auth/database/models"
	"github.com/chrille0234/auth/views/profile"
	"github.com/gin-gonic/gin"
)

func ProfileHandler(ctx *gin.Context) {
	userValue, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userValue.(models.User)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
		return
	}

	handler := templ.Handler(profile.Profile(user))
	handler.ServeHTTP(ctx.Writer, ctx.Request)
}
