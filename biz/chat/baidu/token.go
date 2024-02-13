package baidu

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"runtime/debug"
	"sync"
	"time"
	"web_chat/biz/util/id_gen"
	"web_chat/biz/util/trace_info"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

var (
	accessToken *AppAccessInfo
	mx          sync.Mutex
)

func setAccessToken(info *AppAccessInfo) {
	accessToken = info
}

func getAccessToken(ctx context.Context, appKey, appSecret string) *AppAccessInfo {
	initAccessToken(ctx, appKey, appSecret)
	return accessToken
}

func accessToeknIsNil() bool {
	return accessToken == nil
}

func initAccessToken(ctx context.Context, appKey, appSecret string) error {
	if accessToeknIsNil() {
		mx.Lock()
		defer mx.Unlock()

		if accessToeknIsNil() {
			if err := refreshAccessToken(ctx, appKey, appSecret); err != nil {
				return err
			}

			go func() {
				ctx := context.Background()

				defer func() {
					if rec := recover(); rec != nil {
						hlog.CtxErrorf(ctx, "panic: %v\n %s", rec, debug.Stack())
					}
				}()

				ctx = trace_info.WithTrace(ctx,
					trace_info.TraceInfo{
						LogID: id_gen.NewLogID(),
					})

				timer := time.NewTimer(time.Duration(accessToken.ExpiresIn) * time.Second)
				for range timer.C {
					if err := refreshAccessToken(ctx, appKey, appSecret); err != nil {
						timer.Reset(time.Second)
						continue
					}

					timer.Reset(time.Duration(accessToken.ExpiresIn) * time.Second)
				}
			}()
		}
	}

	return nil
}

func refreshAccessToken(ctx context.Context, appKey, appSecret string) error {
	httpClient := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, accessUrl, nil)
	if err != nil {
		panic(err)
	}

	param := req.URL.Query()
	param.Set("client_id", appKey)
	param.Set("client_secret", appSecret)
	req.URL.RawQuery = param.Encode()

	resp, err := httpClient.Do(req)
	if err != nil {
		hlog.CtxErrorf(ctx, "http request err: %v", err)
		return err
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		hlog.CtxErrorf(ctx, "read resp body err: %v", err)
		return err
	}

	var accessInfo *AppAccessInfo
	if err = json.Unmarshal(data, &accessInfo); err != nil {
		hlog.CtxErrorf(ctx, "unmarshal err: %v", err)
		return err
	}

	setAccessToken(accessInfo)

	return nil
}
