package middleware

import (
	"context"
	"net/http"
	"web_chat/biz/db/redis"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/sessions"
	"github.com/rbcervilla/redisstore/v9"
)

const (
	sessionName        = "session_id"
	sessionStorePrefix = "web_chat_"
)

func SessionMiddleware() app.HandlerFunc {
	store := NewRedisStore()
	store.Options(
		sessions.Options{
			Path:     "/",
			Domain:   "",
			MaxAge:   86400,
			Secure:   false, // https not ready
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
	)
	return sessions.New(sessionName, store)
}

type RedisStore struct {
	*redisstore.RedisStore
}

func (r *RedisStore) Options(opts sessions.Options) {
	r.RedisStore.Options(*opts.ToGorillaOptions())
}

func NewRedisStore() *RedisStore {
	redisStore, err := redisstore.NewRedisStore(context.Background(), redis.GetRedisClient())
	if err != nil {
		panic(err)
	}
	redisStore.KeyPrefix(sessionStorePrefix)
	return &RedisStore{
		RedisStore: redisStore,
	}
}
