package config

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	redisClient "github.com/redis/go-redis/v9"
)

func NewMutex() *redsync.Mutex {
	redis := redisClient.NewClient(&redisClient.Options{
		Addr: "localhost:6379",
	})
	pool := goredis.NewPool(redis)
	rs := redsync.New(pool)
	lock := rs.NewMutex("okegas")
	return lock
}
