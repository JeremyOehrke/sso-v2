package session

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/redis.v3"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func GetBase() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Sup")
	}
}

func SetThing() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var resolvedURL = os.Getenv("REDISCLOUD_URL")
		var password = ""
		if !strings.Contains(resolvedURL, "localhost") {
			parsedURL, _ := url.Parse(resolvedURL)
			password, _ = parsedURL.User.Password()
			resolvedURL = parsedURL.Host
		}

		client := redis.NewClient(&redis.Options{
			Addr:     resolvedURL,
			Password: password,
			DB:       0, // use default DB
		})

		err := client.Set("TEST_KEY", "test_val", 0).Err()
		if err != nil {
			log.Print("error writing key")
		}

		retVal := client.Get("TEST_KEY")
		if retVal.Err() != nil {
			log.Print("error getting key")
		}

		ctx.String(http.StatusOK, "Done! %s -> %s", "TEST_KEY", retVal.Val())
	}
}
