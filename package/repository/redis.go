package repository

import (
	"context"
	// "strconv"
	// "fmt"
	"github.com/go-redis/redis/v8"
)

// var email string

func SetVerificationCode(code string) {

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := rdb.Set(ctx, "verificationCode", code, 0).Err()
	if err != nil {
		panic(err)
	}

}

func GetVerificationCode() string {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	val, err := rdb.Get(ctx, "verificationCode").Result()
	if err != nil {
		panic(err)
	}

	return val

}
