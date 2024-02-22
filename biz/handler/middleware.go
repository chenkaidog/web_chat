package handler

import (
	"bytes"
	"context"
	"net/http"
	"time"
	"web_chat/biz/repository"
	"web_chat/biz/util/origin"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/sessions"
)

func RootMiddleware() []app.HandlerFunc {
	return []app.HandlerFunc{
		BlackListMiddleware(),
	}
}

func AuthMiddleware() []app.HandlerFunc {
	return []app.HandlerFunc{
		LoginRedirectMiddleware(),
		AccountIDMiddleware(),
		AccountStatusMiddleware(),
	}
}

func BlackListMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		isBlock, err := repository.GetRemoteAddrBlackList(ctx, origin.GetIp(c))
		if err != nil {
			c.AbortWithMsg("internal server error", http.StatusInternalServerError)
			return
		}

		if isBlock {
			c.AbortWithMsg("you are blocked", http.StatusForbidden)
			return
		}

		c.Next(ctx)
	}
}

func LoginRedirectMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)

		statusCode := c.Response.StatusCode()
		if statusCode == http.StatusUnauthorized {
			if bytes.Equal(c.Path(), []byte("/login")) {
				return
			}

			c.Redirect(http.StatusOK, []byte("/login"))
		}
	}
}

func AccountIDMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		session := sessions.Default(c)
		accountID, ok := session.Get(sessionAccountID).(string)
		if !ok || accountID == "" {
			c.AbortWithMsg("user not login", http.StatusUnauthorized)
			return
		}
		username, ok := session.Get(sessionAccountUsername).(string)
		if !ok || accountID == "" {
			c.AbortWithMsg("user not login", http.StatusUnauthorized)
			return
		}

		c.Set(sessionAccountID, accountID)
		c.Set(sessionAccountUsername, username)
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
			hlog.CtxErrorf(ctx, "limit usage err: %v", err)
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
