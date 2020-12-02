package main

import (
	"log"
	"os"
	"sso-v2/internal/datasource/redisdatasource"
	"sso-v2/internal/handlers/session"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	/* Dependency Initialization */
	redisUrl := os.Getenv("REDISCLOUD_URL")
	ds := redisdatasource.NewRedisDatasource(redisUrl)
	_ = ds.SetKey("doot", "doot", 0)
	/* End Dependency Initialization */

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	//router.GET("/", session.GetBase())
	router.GET("/sess", session.SetThing())

	router.Run(":" + port)
}
