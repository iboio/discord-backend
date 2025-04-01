package redisdb

import (
	"context"
	"discord/config"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
)

var (
	ctx       = context.Background()
	redisPool = make(map[int]*redis.Client)
	poolLock  sync.Mutex
)

func InitRedis(db int) {
	poolLock.Lock()
	defer poolLock.Unlock()

	if _, exists := redisPool[db]; exists {
		return
	}

	rdb := redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
			Password: "",
			DB:       db,
		})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis connection error (DB %d): %v", db, err)
	}
	redisPool[db] = rdb
	fmt.Printf("Connected to Redis (DB %d)!\n", db)
}

func GetRedisClient(db int) *redis.Client {
	poolLock.Lock()
	defer poolLock.Unlock()

	if _, exists := redisPool[db]; !exists {
		InitRedis(db)
	}

	return redisPool[db]
}
