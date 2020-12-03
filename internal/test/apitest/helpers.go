package apitest

import "github.com/gin-gonic/gin"

func BuildTestRouter(method string, route string, handlerFunc gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Handle(method, route, handlerFunc)
	return r
}
