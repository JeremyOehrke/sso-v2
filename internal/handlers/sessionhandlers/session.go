package sessionhandlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sso-v2/internal/handlers"
	"sso-v2/internal/service/session"
)

func GetSessionDataHandler(svc session.SessionSVC) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId := ctx.Param("sessionId")
		if sessionId == "" {
			ctx.Data(http.StatusNotFound, gin.MIMEPlain, nil)
			return
		}

		sessionData, err := svc.GetSessionById(sessionId)
		if err == session.SessionNotFoundError { //if the session isn't found, don't log
			ctx.Data(http.StatusNotFound, gin.MIMEPlain, nil)
			return
		}
		if err != nil { //log and return for all others
			log.Printf("error looking up session: %v", err.Error())
			ctx.JSON(http.StatusInternalServerError, handlers.ErrorMessage{Message: "error locating session"})
			return
		}

		ctx.JSON(http.StatusOK, *sessionData)
	}
}

type setSessionRequest struct {
	SessionVars map[string]string `json:"sessionVars"`
}

func SetSessionDataHandler(svc session.SessionSVC) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId := ctx.Param("sessionId")
		if sessionId == "" {
			ctx.Data(http.StatusNotFound, gin.MIMEPlain, nil)
			return
		}

		requestData := &setSessionRequest{}
		err := ctx.BindJSON(requestData)
		if err != nil {
			log.Printf("error binding request body: %v", err.Error())
			ctx.JSON(http.StatusInternalServerError, handlers.ErrorMessage{Message: "error processing body"})
			return
		}
		if requestData.SessionVars == nil {
			ctx.JSON(http.StatusBadRequest, handlers.ErrorMessage{Message: "missing body"})
			return
		}

		fmt.Println(sessionId)
		err = svc.SetSessionBodyById(sessionId, requestData.SessionVars)
		if err == session.SessionNotFoundError {
			ctx.Data(http.StatusNotFound, gin.MIMEPlain, nil)
		}
		ctx.Data(http.StatusNoContent, gin.MIMEPlain, nil)
	}
}
