package routes

import (
	"github.com/gin-gonic/gin"
	"sso-v2/internal/handlers/userhandlers"
	"sso-v2/internal/service/session"
	"sso-v2/internal/service/user"
)

func BuildRoutes(ginMode string, usersvc user.UserSVC, sessionsvc session.SessionSVC) *gin.Engine {
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
			usrs.POST("/doAuth", userhandlers.AuthUserHandler(usersvc))
		}

		sess := v1.Group("/sessions")
		{
			sess.POST("/", func(context *gin.Context) {

			})
		}
	}

	return router
}
