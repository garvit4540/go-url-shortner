package database

import (
	"context"
	"github.com/garvit4540/go-url-shortner/trace"
	redis "github.com/go-redis/redis/v8"
	"os"
)

var Ctx = context.Background()

func CreateClient(dbNo int) *redis.Client {

	// Initialising a redis db
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDR"), // Redis server address
		Password: os.Getenv("DB_PASS"), // No password set
		DB:       dbNo,                 // Default DB
	})

	// Test the connection with a Ping
	err := rdb.Ping(Ctx).Err()
	if err != nil {
		trace.LogError(trace.ErrorConnectingToRedis, err, map[string]interface{}{"Database No. ": dbNo})
	}
	trace.LogInfo(trace.SuccessfullyConnectedToRedis, map[string]interface{}{"Database No. ": dbNo})

	return rdb

}
