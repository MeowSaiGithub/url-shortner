package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
)

func CreateClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST"),
	})
	err := rdb.Ping(context.Background())
	return rdb, err.Err()
}
