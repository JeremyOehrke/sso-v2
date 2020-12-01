package session

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
)

func GetBase() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Sup")
	}
}

func SetThing() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		rdb := redis.NewClient(&redis.Options{
			Addr:     "redis-12250.c11.us-east-1-2.ec2.cloud.redislabs.com:12250",
			Password: "Vfzl3XZbh0A7vcIl8ZZsOgbn2XkFuGu6", // no password set
			DB:       0,                                  // use default DB
		})

		err := rdb.Set(ctx, "TEST_KEY", "SomeValue", 0).Err()
		if err != nil {
			log.Print("Error setting key")
		}

		val, err := rdb.Get(ctx, "TEST_KEY").Result()
		if err != nil {
			log.Print("Error getting key")
		}

		ctx.String(http.StatusOK, "Done! %s -> %s", "TEST_KEY", val)
	}
}
