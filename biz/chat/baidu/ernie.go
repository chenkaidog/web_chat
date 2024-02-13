package baidu

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"web_chat/biz/config"
	"web_chat/biz/model/domain"
	"web_chat/biz/util/sse_client"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func NewErine(model string) *Erine {
	switch model {
	case domain.ModelErine4:
		return &Erine{
			ChatUrl:   erine4chatUrl,
			AppKey:    config.GetBaiduConf().AppKey,
			AppSecret: config.GetBaiduConf().AppSecret,
		}
	}

	return nil
}

type Erine struct {
	AppKey    string
	AppSecret string
	ChatUrl   string
}

func (e *Erine) StreamChat(ctx context.Context, chatContext []*domain.ChatContent) (chan string, chan *domain.PlatformError, error) {
	httpReq, err := e.newStreamChatRequest(ctx, chatContext)
	if err != nil {
		return nil, nil, err
	}

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		hlog.CtxErrorf(ctx, "http request err: %v", err)
		return nil, nil, err
	}

	if httpResp.StatusCode != http.StatusOK {
		respContent, err := io.ReadAll(httpResp.Body)
		defer httpResp.Body.Close()
		if err != nil {
			return nil, nil, err
		}

		hlog.CtxErrorf(ctx, "status_code: %d, error_msg: %s", httpResp.StatusCode, respContent)
		return nil, nil, errors.New("request fails")
	}

	errCh := make(chan *domain.PlatformError)
	respch := sse_client.HandleSseResp(ctx, httpResp, func(ctx context.Context, data []byte) (string, bool) {
		var respBody ChatCreateResp
		if err := json.Unmarshal(data, &respBody); err != nil {
			hlog.CtxErrorf(ctx, "json unmarshal err: %v", err)
			return "", false
		}

		if respBody.ErrorCode != 0 || respBody.Error != "" {
			hlog.CtxErrorf(ctx, "request err: %s. %s.", respBody.ErrorMsg, respBody.ErrorDescription)
			errCh <- &domain.PlatformError{
				Code: respBody.ErrorCode,
				Msg:  respBody.ErrorMsg,
			}
			return "", true
		}

		if respBody.IsEnd {
			return "", true
		}

		if len(respBody.Result) > 0 {
			return respBody.Result, false
		}

		return "", false
	})

	return respch, errCh, nil
}

func (e *Erine) newStreamChatRequest(ctx context.Context, chatContext []*domain.ChatContent) (*http.Request, error) {
	var messages []*Message
	for _, content := range chatContext {
		messages = append(
			messages,
			&Message{
				Role:    roleMapper[content.Role],
				Content: content.Content,
			},
		)
	}

	reqBody, err := json.Marshal(
		&ChatCreateReq{
			Stream:   true,
			Messages: messages,
		},
	)
	if err != nil {
		hlog.CtxErrorf(ctx, "json marshal err: %v", err)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, e.ChatUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		hlog.CtxErrorf(ctx, "new request err: %v", err)
		return nil, err
	}

	param := req.URL.Query()
	param.Set("access_token", getAccessToken(ctx, e.AppKey, e.AppSecret).AccessToken)
	req.URL.RawQuery = param.Encode()

	return req, nil
}
