package configs

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	ErrLoadRedisEnvFile    = errors.New("load .env file failure")
	ErrConnect2RedisServer = errors.New("connect to redis server via Redis Client failure")
)

type redisCacheClient struct {
	instance *redis.Client
}

func NewRedisCacheClient() *redisCacheClient {
	return &redisCacheClient{instance: nil}
}

func (instance *redisCacheClient) Instance() (*redis.Client, error) {
	if instance.instance == nil {
		if err := godotenv.Load(); err != nil {
			fmt.Println("Error while load .env file: " + err.Error())
			return nil, ErrLoadGormEnvFile
		} else {
			var dns = os.Getenv("REDIS_URL")
			if opt, err := redis.ParseURL(dns); err != nil {
				fmt.Println("Can't connect to Redis Server: ", err.Error())
				return nil, ErrConnect2RedisServer
			} else {
				instance.instance = redis.NewClient(opt)
			}
		}
	}
	return instance.instance, nil
}
