package redis

import (
	"aapi/config"
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

// Return new New Redis db instance
func NewRedis(c *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.RedisAddr,
		Password: c.Redis.RedisPassword, // no password set
		DB:       0,                     // use default DB
		PoolSize: c.Redis.PoolSize,
	})

	return rdb, nil
}

// Return new New Redis replicas db instance
func NewReplicaRedis(c *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.RedisReplica.RedisAddr,
		Password: c.RedisReplica.RedisPassword, // no password set
		DB:       0,                            // use default DB
		PoolSize: c.RedisReplica.PoolSize,
	})

	return rdb, nil
}

func GetValue(ctx context.Context, rd *redis.Client, key string) (string, error) {
	cmd := rd.Get(ctx, key)
	if cmd == nil {
		log.Default().Printf("Key (%v) Not found", key)
		return "", errors.New("redis does not work")
	}

	val, err := cmd.Result()
	if err != nil || err == redis.Nil {
		log.Default().Printf("Key (%v) Empty result", key)
		return "", nil
	}
	return val, err
}

func GetUnmarshallValue(ctx context.Context, rd *redis.Client, key string, output interface{}) error {
	jsonString, err := GetValue(ctx, rd, key)
	if err != nil {
		return err
	}

	if jsonString == "" {
		return redis.Nil
	}

	return json.Unmarshal([]byte(jsonString), &output)
}

func SetStringValue(ctx context.Context, rd *redis.Client, key string, val interface{}, ttl time.Duration) error {
	log.Default().Printf("SetStringValue (%v) ", key)
	switch val.(type) {
	case int, int32, int64, float32, float64, string:
		return rd.Set(ctx, key, val, ttl).Err()
	default:
		jsonBuf, err := json.Marshal(val)

		if err != nil {
			return err
		}

		return rd.Set(ctx, key, jsonBuf, ttl).Err()
	}
}

func RemoveCacheByKey(ctx context.Context, rd *redis.Client, key string) error {
	log.Default().Printf("RemoveCacheByKey (%v) ", key)
	return rd.Del(ctx, key).Err()
}
