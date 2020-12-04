package routes

import (
	"github.com/gin-gonic/gin"
	"sso-v2/internal/handlers/sessionhandlers"
	"sso-v2/internal/handlers/userhandlers"
	"sso-v2/internal/service/session"
	"sso-v2/internal/service/user"
)

func BuildRouter(ginMode string, usersvc user.UserSVC, sessionsvc session.SessionSVC) *gin.Engine {
	gin.SetMode(ginMode)
	router := gin.Default()
	router.Use(gin.Logger())

	//V1 routes
	v1 := router.Group("/v1")
	{
		//User Routes
		usrs := v1.Group("/users")
		{
			usrs.POST("/", userhandlers.CreateUserHandler(usersvc))
			usrs.POST("/doAuth", userhandlers.AuthUserHandler(usersvc, sessionsvc))
		}
		//Session routes
		sess := v1.Group("/sessions")
		{
			sess.GET("/:sessionId", sessionhandlers.GetSessionDataHandler(sessionsvc))
			sess.PUT("/:sessionId", sessionhandlers.SetSessionDataHandler(sessionsvc))
			sess.DELETE("/:sessionId", sessionhandlers.DestroySessionHandler(sessionsvc))
		}
	}

	return router
}
