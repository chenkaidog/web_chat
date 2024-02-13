package sse_client

import (
	"bufio"
	"bytes"
	"context"
	"net/http"
	"runtime/debug"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type Handler func(ctx context.Context, msg []byte) (string, bool)

func HandleSseResp(ctx context.Context, resp *http.Response, handle Handler) chan string {
	eventCh := readStreamResp(ctx, resp)
	respCh := make(chan string)

	go func() {
		defer func() {
			close(respCh)
			if rec := recover(); rec != nil {
				hlog.CtxErrorf(ctx, "panic recover: rec: %v\n%s", rec, debug.Stack())
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-eventCh:
				if !ok {
					return
				}
				if event != nil && event.Data != nil {
					hlog.CtxDebugf(ctx, "msg: %s", event.Data)
					result, stop := handle(ctx, event.Data)
					if stop {
						return
					}
					if result != "" {
						respCh <- result
					}
				}
			}
		}
	}()

	return respCh
}

func readStreamResp(ctx context.Context, resp *http.Response) chan *SseEvent {
	eventCh := make(chan *SseEvent)

	go func() {
		defer func() {
			close(eventCh)
			resp.Body.Close()

			if rec := recover(); rec != nil {
				hlog.CtxErrorf(ctx, "panic recover: rec: %v\n%s", rec, debug.Stack())
			}
		}()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				hlog.CtxErrorf(ctx, "scanner err: %v", err)
				return
			}

			data := scanner.Bytes()
			hlog.CtxInfof(ctx, "msg: %s", data)

			select {
			case eventCh <- parseEvent(data):
			case <-ctx.Done():
				return
			}
		}
	}()

	return eventCh
}

type SseEvent struct {
	ID    []byte
	Data  []byte
	Event []byte
	Retry []byte
}

var (
	headerID    = []byte("id:")
	headerData  = []byte("data:")
	headerEvent = []byte("event:")
	headerRetry = []byte("retry:")
)

func parseEvent(msg []byte) *SseEvent {
	if len(msg) <= 0 {
		return nil
	}

	var e SseEvent

	// Normalize the crlf to lf to make it easier to split the lines.
	// Split the line by "\n" or "\r", per the spec.
	for _, line := range bytes.FieldsFunc(msg, func(r rune) bool { return r == '\n' || r == '\r' }) {
		switch {
		case bytes.HasPrefix(line, headerID):
			e.ID = append([]byte(nil), trimHeader(len(headerID), line)...)
		case bytes.HasPrefix(line, headerData):
			// The spec allows for multiple data fields per event, concatenated them with "\n".
			e.Data = append(e.Data[:], append(trimHeader(len(headerData), line), byte('\n'))...)
		// The spec says that a line that simply contains the string "data" should be treated as a data field with an empty body.
		case bytes.Equal(line, bytes.TrimSuffix(headerData, []byte(":"))):
			e.Data = append(e.Data, byte('\n'))
		case bytes.HasPrefix(line, headerEvent):
			e.Event = append([]byte(nil), trimHeader(len(headerEvent), line)...)
		case bytes.HasPrefix(line, headerRetry):
			e.Retry = append([]byte(nil), trimHeader(len(headerRetry), line)...)
		default:
			// Ignore any garbage that doesn't match what we're looking for.
		}
	}

	// Trim the last "\n" per the spec.
	e.Data = bytes.TrimSuffix(e.Data, []byte("\n"))

	return &e
}

func trimHeader(size int, data []byte) []byte {
	if data == nil || len(data) < size {
		return data
	}

	data = data[size:]
	// Remove optional leading whitespace
	if len(data) > 0 && data[0] == 32 {
		data = data[1:]
	}
	// Remove trailing new line
	if len(data) > 0 && data[len(data)-1] == 10 {
		data = data[:len(data)-1]
	}
	return data
}
