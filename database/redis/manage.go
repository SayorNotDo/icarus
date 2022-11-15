package redismanage

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

var RedisString = redis.String

var Pool *redis.Pool

func init() {
	Pool = &redis.Pool{
		MaxIdle:     20,
		MaxActive:   0,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
	log.Println("--------------------------redis connection pool--------------------------")
}
