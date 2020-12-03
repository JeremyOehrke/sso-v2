package userhandlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sso-v2/internal/handlers"

	"sso-v2/internal/service/user"
)

type createUserBody struct {
	Username string
	Password string
}

func CreateUserHandler(svc user.UserSVC) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userData := &createUserBody{}

		err := ctx.BindJSON(userData)
		if err != nil {
			log.Printf("error binding request body: %v", err.Error())
			ctx.JSON(http.StatusInternalServerError, handlers.ErrorMessage{Message: "error processing body"})
			return
		}
		if userData.Username == "" || userData.Password == "" {
			ctx.JSON(http.StatusBadRequest, handlers.ErrorMessage{Message: "missing username and/or password"})
			return
		}

		err = svc.CreateUser(userData.Username, userData.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, handlers.ErrorMessage{Message: "error creating user"})
			return
		}
		ctx.Data(http.StatusCreated, gin.MIMEPlain, nil)
	}
}
