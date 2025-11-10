package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()

// accountsCacheTTL returns TTL for accounts cache from env REDIS_ACCOUNTS_TTL or defaults to 15m
func accountsCacheTTL() time.Duration {
	s := os.Getenv("REDIS_ACCOUNTS_TTL")
	if s == "" {
		return 15 * time.Minute
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 15 * time.Minute
	}
	return d
}

func InitRedisConnection() error {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found, using default Redis settings")
	}

	// Получаем настройки из переменных окружения или используем дефолтные
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	password := os.Getenv("REDIS_PASSWORD")

	dbStr := os.Getenv("REDIS_DB")
	db := 0
	if dbStr != "" {
		if parsedDB, err := strconv.Atoi(dbStr); err == nil {
			db = parsedDB
		}
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	// Тестируем подключение
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis connection test failed: %v", err)
		return err
	}

	log.Printf("Redis connected successfully to %s:%s", host, port)
	return nil
}

func GetRedisClient() *redis.Client {
	return rdb
}

// SetAccountsCache - кеширует счета пользователя на 15 минут
func SetAccountsCache(userID int, accounts interface{}) error {
	if rdb == nil {
		return fmt.Errorf("redis client not initialized")
	}

	key := fmt.Sprintf("user_accounts:%d", userID)

	data, err := json.Marshal(accounts)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, data, accountsCacheTTL()).Err()
}

// GetAccountsCache - получает кешированные счета пользователя
func GetAccountsCache(userID int, result interface{}) error {
	if rdb == nil {
		return fmt.Errorf("redis client not initialized")
	}

	key := fmt.Sprintf("user_accounts:%d", userID)

	data, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), result)
}

// DeleteAccountsCache - удаляет кеш счетов пользователя
func DeleteAccountsCache(userID int) error {
	if rdb == nil {
		return nil // Не возвращаем ошибку если Redis не инициализирован
	}

	key := fmt.Sprintf("user_accounts:%d", userID)
	return rdb.Del(ctx, key).Err()
}

// DeleteAccountCacheByAccountID - удаляет кеш по ID аккаунта (нужно получить userID)
func DeleteAccountCacheByAccountID(accountID int) error {
	if rdb == nil {
		return nil // Не возвращаем ошибку если Redis не инициализирован
	}

	// Можно добавить дополнительную логику для получения userID по accountID
	// Для простоты пока удаляем все кеши аккаунтов
	pattern := "user_accounts:*"
	keys, err := rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return rdb.Del(ctx, keys...).Err()
	}

	return nil
}

// Ping checks Redis connectivity
func Ping() error {
	if rdb == nil {
		return fmt.Errorf("redis client not initialized")
	}
	_, err := rdb.Ping(ctx).Result()
	return err
}
