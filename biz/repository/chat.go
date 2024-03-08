package repository

import (
	"context"
	"fmt"
	"time"
	"web_chat/biz/db/redis"
)

const qpsLimit = time.Second

func QPSLimitBySession(ctx context.Context, sessID,path string) (bool, error) {
	return redis.GetRedisClient().
		SetNX(ctx, genUsageLimitKey(sessID,path), true, qpsLimit).
		Result()
}

func genUsageLimitKey(sessID,path string) string {
	return fmt.Sprintf("usage_limit_key_%s_%s", sessID,path)
}
