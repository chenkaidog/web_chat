package handler

import (
	"context"
	"net/http"
	"time"
	"web_chat/biz/repository"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/sessions"
)

func AccountIDMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		session := sessions.Default(c)
		if session == nil {
			c.AbortWithMsg("user not login", http.StatusUnauthorized)
			return
		}
		accountID, ok := session.Get(sessionAccountID).(string)
		if !ok || accountID == "" {
			c.AbortWithMsg("user not login", http.StatusUnauthorized)
			return
		}
		c.Set(sessionAccountID, accountID)
		c.Set(sessionSessID, session.ID())

		c.Next(ctx)
	}
}

func AccountStatusMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		accountID := c.GetString(sessionAccountID)
		sessID := c.GetString(sessionSessID)

		sessionList, err := repository.GetSessionList(ctx, accountID)
		if err != nil {
			hlog.CtxErrorf(ctx, "session list err: %v", err)
			c.AbortWithMsg("internal server error", http.StatusInternalServerError)
			return
		}
		if sessionList == nil || !containSession(sessionList, sessID) {
			c.AbortWithMsg("user not login", http.StatusUnauthorized)
			return
		}

		account, err := repository.GetAccountByAccountID(ctx, accountID)
		if err != nil || account == nil {
			hlog.CtxErrorf(ctx, "account is empty %v, or err: %v", account, err)
			c.AbortWithMsg("internal server error", http.StatusInternalServerError)
			return
		}

		if account.IsInvalid() {
			c.AbortWithMsg("account status is invalid", http.StatusForbidden)
			return
		}

		if account.ExpirationDate.Before(time.Now()) {
			c.AbortWithMsg("account is expired", http.StatusForbidden)
			return
		}

		c.Next(ctx)
	}
}

func containSession(sessionList []string, sessID string) bool {
	for _, id := range sessionList {
		if id == sessID {
			return true
		}
	}

	return false
}

func ChatLimitMiddleware() []app.HandlerFunc {
	return []app.HandlerFunc{
		chatQpsLimitMiddleware(),
	}
}

func chatQpsLimitMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		sessID := c.GetString(sessionSessID)
		ok, err := repository.QPSLimitBySession(ctx, sessID)
		if err != nil {
			hlog.CtxErrorf(ctx, "limit usage err err: %v", err)
			c.AbortWithMsg("internal server error", http.StatusInternalServerError)
			return
		}

		if !ok {
			c.AbortWithMsg("qps too high", http.StatusForbidden)
			return
		}

		c.Next(ctx)
	}
}
