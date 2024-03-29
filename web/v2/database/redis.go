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

	keys, err := r.Keys(ctx, keysPattern).Result()
	if err != nil {
		log.Println("Error geting cache keys to clear:", err)
		cacheStatus = "stale"
	}

	if keys != nil && len(keys) > 0 {
		err := r.Del(ctx, keys...).Err()
		if err != nil {
			log.Println("Error clearing cache keys:", err)
			cacheStatus = "stale"
		}
	}

	return cacheStatus
}

func (r *Redis) GetDecodedCache(ctx context.Context, userID int, namespace string, dest interface{}) error {
	key := fmt.Sprintf("db:%v:the-route-doesnt-matter-here:%v", userID, namespace)

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

func (r *Redis) SetEncodedCache(ctx context.Context, userID int, namespace string, value interface{}, exp time.Duration) error {
	key := fmt.Sprintf("db:%v:the-route-doesnt-matter-here:%v", userID, namespace)

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

func (r *Redis) PurgeAllCachedData(ctx context.Context) string {
	fmt.Println("Purging all cached data, but preserving user sessions...")

	keysPattern := "*:*:*:*"
	cacheStatus := "cleared"

	keys, err := r.Keys(ctx, keysPattern).Result()
	if err != nil {
		log.Println("Error geting cache keys to clear:", err)
		cacheStatus = "stale"
	}

	if keys != nil && len(keys) > 0 {
		err := r.Del(ctx, keys...).Err()
		if err != nil {
			log.Println("Error clearing cache keys:", err)
			cacheStatus = "stale"
		}
	}

	fmt.Println("Cache status is:", cacheStatus)

	return cacheStatus
}
