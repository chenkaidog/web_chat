package middleware

import (
	"context"
	"web_chat/biz/util/id_gen"
	"web_chat/biz/util/trace_info"

	"github.com/cloudwego/hertz/pkg/app"
)

const (
	headerKeyTraceId = "X-Trace-ID"
	headerKeyLogId   = "X-Log-ID"
	headerKeySpanId  = "X-Span-ID"
)

func TraceContextMW() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		logID := c.Request.Header.Get(headerKeyLogId)
		if logID == "" {
			logID = id_gen.NewLogID()
		}

		ctx = trace_info.WithTrace(
			ctx,
			trace_info.TraceInfo{
				LogID: logID,
			})

		c.Next(ctx)

		c.Header(headerKeyLogId, logID)
	}
}
