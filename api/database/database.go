package database

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func CreateClient(dbbase int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDR"),
		Password: os.Getenv("DB_PASS"),
		DB:       dbbase,
	})

	return rdb
}
