package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"os"
	"sso-v2/internal/datasource/redisdatasource"
	"sso-v2/internal/handlers/routes"
	"sso-v2/internal/service/session/sessionsvc"
	"sso-v2/internal/service/user/usersvc"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	/* Dependency Initialization */
	redisUrl := os.Getenv("REDISCLOUD_URL")
	ds := redisdatasource.NewRedisDatasource(redisUrl)
	userSvc := usersvc.NewUserSvc(ds)
	sessionSvc := sessionsvc.NewSessionSvc(ds)
	/* End Dependency Initialization */

	router := routes.BuildRoutes(gin.ReleaseMode, userSvc, sessionSvc)
	router.Use(gin.Logger())

	router.Run(":" + port)
}
