package middleware

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func recoveryHandler(c context.Context, ctx *app.RequestContext, err interface{}, stack []byte) {
	hlog.CtxErrorf(c, "[Recovery] err=%v\nstack=%s", err, stack)
	ctx.AbortWithStatus(http.StatusInternalServerError)
}

func RecoveryMiddleware() app.HandlerFunc {
	return recovery.Recovery(recovery.WithRecoveryHandler(recoveryHandler))
}
