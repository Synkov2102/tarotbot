package redisConnector

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
	once        sync.Once
	ctx         = context.Background()
)

// GetRedisClient возвращает экземпляр клиента Redis, используя шаблон одиночка.
func GetRedisClient() *redis.Client {
	once.Do(func() {
		redisAddr := os.Getenv("REDIS_HOST") // Замените на ваш адрес Redis
		redisPassword := ""                  // Замените на ваш пароль, если есть

		redisClient = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       0, // Используемая база данных
		})

		_, err := redisClient.Ping(ctx).Result()
		if err != nil {
			log.Fatalf("Ошибка подключения к Redis: %v", err)
		}
	})
	return redisClient
}
