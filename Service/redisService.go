package service

import (
	datamodel "ToDoApp/DataModel"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func InitRedisClient() {
	redisClient = ConnectToRedis("localhost:6379", "", 0)
}

func ConnectToRedis(url, password string, DB int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "",
		DB:       DB,
	})
}

func CheckInRedis(sessionId string) bool {
	val, err := redisClient.Get(sessionId).Result()
	if err != nil {
		fmt.Println(err)
	}
	if len(val) != 0 {
		return true
	}
	return false
}

func storeSession(sessionId string, user datamodel.RedisUser) {
	fmt.Println("********************>>", sessionId)
	json, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	err = redisClient.Set(sessionId, json, 0).Err()
	if err != nil {
		panic(err)
	}
}
