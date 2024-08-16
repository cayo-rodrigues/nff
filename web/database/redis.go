package database

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	*redis.Client
}

func (r *Redis) ClearCache(ctx context.Context, userID int, namespace string) string {
	keysPattern := fmt.Sprintf("*:%v:*:%v", userID, namespace)
	cacheStatus := "cleared"

	err := r.DelMany(ctx, keysPattern)
	if err != nil {
		cacheStatus = "stale"
	}

	return cacheStatus
}

func (r *Redis) DestroyAllCachedData(ctx context.Context) string {
	fmt.Println("Destroying all cached data, but preserving user sessions...")

	keysPattern := "*:*:*:*"
	cacheStatus := "cleared"

	err := r.DelMany(ctx, keysPattern)
	if err != nil {
		cacheStatus = "stale"
	}

	fmt.Println("Cache status is:", cacheStatus)

	return cacheStatus
}

func (r *Redis) GetDecodedCache(ctx context.Context, userID int, namespace string, filters string, dest interface{}) error {
	key := fmt.Sprintf("db:%v:%v:%v", userID, filters, namespace)

	cachedValue, err := r.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return err
	}
	if err != nil {
		log.Printf("Error trying to get %s cache at db level: %v\n", namespace, err)
		return err
	}

	decoder := gob.NewDecoder(bytes.NewReader(cachedValue))
	err = decoder.Decode(dest)
	if err != nil {
		log.Printf("Error trying to decode %s cache at db level: %v\n", namespace, err)
		return err
	}

	return nil
}

func (r *Redis) SetEncodedCache(ctx context.Context, userID int, namespace string, filters string, value interface{}, exp time.Duration) error {
	key := fmt.Sprintf("db:%v:%v:%v", userID, filters, namespace)

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(value)
	if err != nil {
		log.Printf("Error trying to encode %s cache at db level: %v\n", namespace, err)
		return err
	}

	err = r.Set(ctx, key, buf.Bytes(), exp).Err()
	if err != nil {
		log.Printf("Error trying to set %s cache at db level: %v\n", namespace, err)
		return err
	}

	return nil
}

func (r *Redis) DelMany(ctx context.Context, keysPattern string) error {
	var cursor uint64
	var batchSize int64 = 100

	scanner := r.Scan(ctx, cursor, keysPattern, batchSize).Iterator()

	keys := []string{}

	for scanner.Next(ctx) {
		keys = append(keys, scanner.Val())
	}

	err := scanner.Err()
	if err != nil {
		log.Println("Error scanning cache keys:", err)
		return err
	}

	if len(keys) > 0 {
		err = r.Del(ctx, keys...).Err()
		if err != nil {
			log.Println("Error clearing cache keys:", keys, err)
			return err
		}
	}

	return nil
}
