package repository

import (
	"context"
	"fmt"
	"time"
	"web_chat/biz/db/redis"
)

const qpsLimit = 5 * time.Second

func QPSLimitBySession(ctx context.Context, sessID string) (bool, error) {
	return redis.GetRedisClient().
		SetNX(ctx, genUsageLimitKey(sessID), true, qpsLimit).
		Result()
}

func genUsageLimitKey(sessID string) string {
	return fmt.Sprintf("usage_limit_key_%s", sessID)
}
