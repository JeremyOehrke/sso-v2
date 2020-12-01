package session

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetBase() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.String(http.StatusOK, "Sup")
	}
}