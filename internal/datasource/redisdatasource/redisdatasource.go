package redisdatasource

import (
	"gopkg.in/redis.v3"
	"log"
	"net/url"
	"sso-v2/internal/datasource"
	"strings"
	"time"
)

type RedisDataSource struct {
	cli *redis.Client
}

func NewRedisDatasource(dsUrl string) datasource.Datasource {
	rds := &RedisDataSource{}
	password := ""
	resolvedURL := ""
	if !strings.Contains(dsUrl, "localhost") {
		parsedURL, _ := url.Parse(dsUrl)
		password, _ = parsedURL.User.Password()
		resolvedURL = parsedURL.Host
	}

	rds.cli = redis.NewClient(&redis.Options{
		Addr:     resolvedURL,
		Password: password,
		DB:       0, // use default DB
	})
	return rds
}

func (ds *RedisDataSource) GetKey(key string) (val string, err error) {
	retVal, err := ds.cli.Get(key).Result()

	// This reset of the returned error is to prevent having to handle the no record found error across the paplication
	// since this application doesn't view no record to be an error
	if err == redis.Nil {
		err = nil
	}

	if err != nil {
		log.Print("error getting key: " + err.Error())
	}
	return retVal, err
}

func (ds *RedisDataSource) SetKey(key string, val string, timeout time.Duration) error {
	err := ds.cli.Set(key, val, timeout).Err()
	if err != nil {
		log.Print("error writing key: " + err.Error())
	}
	return err
}

func (ds *RedisDataSource) DelKey(key string) error {
	err := ds.cli.Del(key).Err()
	if err != nil {
		log.Print("error deleting key: " + err.Error())
	}
	return err
}
