package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

func FlushRedisCache(redisClient *redis.Client, trigger int, log *logger.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch trigger {
	case 0:
		log.Info("🔒 Redis cache preservation mode")
		return nil

	case 1:
		log.Info("🧹 Performing selective Redis cache cleanup...")

		// То что мы будем удалять
		patterns := []string{
			"user:*",
			"session:*",
			"cache:*",
		}

		var deletedKeysCount int64
		for _, pattern := range patterns {
			keys, err := redisClient.Keys(ctx, pattern).Result()
			if err != nil {
				log.Errorf("❌ Error finding keys for pattern %s: %v", pattern, err)
				continue
			}

			if len(keys) > 0 {
				deletedCount, err := redisClient.Del(ctx, keys...).Result()
				if err != nil {
					log.Errorf("❌ Error deleting keys for pattern %s: %v", pattern, err)
				}
				deletedKeysCount += deletedCount
			}
		}

		log.Infof("✅ Selective cache cleanup complete. Deleted %d keys", deletedKeysCount)
		return nil

	case 2:
		log.Warn("🚨 Performing FULL Redis cache destruction...")

		err := redisClient.FlushAll(ctx).Err()
		if err != nil {
			log.Errorf("❌ Full Redis cache flush failed: %v", err)
			return err
		}

		log.Info("💥 FULL Redis cache successfully destroyed")
		return nil

	default:
		log.Warn("❓ Invalid trigger value. No action taken.")
		return fmt.Errorf("invalid trigger value")
	}
}
