package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/cayo-rodrigues/nff/web/internal/db"
)

func ClearCache(ctx context.Context, userID int, namespace string) string {
	keysPattern := fmt.Sprintf("*:%v:*:%v", userID, namespace)
	cacheStatus := "cleared"

	keys, err := db.Redis.Keys(ctx, keysPattern).Result()
	if err != nil {
		log.Println("Error geting cache keys to clear:", err)
		cacheStatus = "stale"
	}

	if keys != nil && len(keys) > 0 {
		err := db.Redis.Del(ctx, keys...).Err()
		if err != nil {
			log.Println("Error clearing cache keys:", err)
			cacheStatus = "stale"
		}
	}

	return cacheStatus
}

