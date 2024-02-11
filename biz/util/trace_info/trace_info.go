package trace_info

import (
	"context"
)

type TraceInfo struct {
	LogID string
}

func WithTrace(ctx context.Context, trace TraceInfo) context.Context {
	ctx = context.WithValue(
		ctx,
		TraceInfo{},
		trace,
	)

	return ctx
}

func GetTraceInfo(ctx context.Context) TraceInfo {
	if ctx == nil {
		return TraceInfo{}
	}

	trace, ok := ctx.Value(TraceInfo{}).(TraceInfo)
	if ok {
		return trace
	}

	return TraceInfo{}
}
