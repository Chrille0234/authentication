package main

import (
	"github.com/chrille0234/auth/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load environment variables, so they can be accessed with os.Getenv
	godotenv.Load()

	db := database.ConnectAndSeed("db.sqlite3")
	defer db.Close()

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.LoadHTMLGlob("./static/*.css")
	router.Static("/static", "./static")

    SetupAllRoutes(router, db)

	router.Run(":4000")
}
