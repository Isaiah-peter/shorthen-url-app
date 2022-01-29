package database

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

 var ctx = context.Background()
 func CreateClient(dbbase int) *redis.Client{
	 rdb := redis.NewClient(&redis.Options{
		 Addr: os.Getenv("") ,
		 Password: ,
		 DB: dbdbbase,
	 })
 }