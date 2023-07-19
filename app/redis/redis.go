package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
)

var ctx = context.Background()

func CreateClient(dbNo int) *redis.Client {
	clientOptions := &redis.Options{
		Addr:     os.Getenv("DB_ADDR"),
		Password: os.Getenv("DB_PASS"),
		DB:       dbNo,
	}

	return redis.NewClient(clientOptions)
}
