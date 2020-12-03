package userhandlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sso-v2/internal/handlers"

	"sso-v2/internal/service/user"
)

type userRequestBody struct {
	Username string
	Password string
}

func CreateUserHandler(svc user.UserSVC) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userData, done := bindRequestData(ctx)
		if done {
			return
		}

		err := svc.CreateUser(userData.Username, userData.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, handlers.ErrorMessage{Message: "error creating user"})
			return
		}
		ctx.Data(http.StatusCreated, gin.MIMEPlain, nil)
	}
}

type authResponse struct {
	AuthOk bool `json:"authOk"`
}

func AuthUserHandler(svc user.UserSVC) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userData, done := bindRequestData(ctx)
		if done {
			return
		}

		authed, err := svc.AuthUser(userData.Username, userData.Password)
		//This only logs and sends an error if we got some other error than the user just not being found
		//User not found is an expected and acceptable edge case we wouldn't want to page on
		if err != nil && err != user.NotFound {
			log.Printf("error authorizing user: %v", err.Error())
			ctx.JSON(http.StatusInternalServerError, handlers.ErrorMessage{Message: "error authorizing user"})
			return
		}

		ctx.JSON(http.StatusOK, authResponse{AuthOk: authed})
	}
}

func bindRequestData(ctx *gin.Context) (*userRequestBody, bool) {
	userData := &userRequestBody{}

	err := ctx.BindJSON(userData)
	if err != nil {
		log.Printf("error binding request body: %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, handlers.ErrorMessage{Message: "error processing body"})
		return nil, true
	}
	if userData.Username == "" || userData.Password == "" {
		ctx.JSON(http.StatusBadRequest, handlers.ErrorMessage{Message: "missing username and/or password"})
		return nil, true
	}
	return userData, false
}
