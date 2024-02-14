package baidu

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"web_chat/biz/config"
	"web_chat/biz/model/domain"
	"web_chat/biz/model/err"
	"web_chat/biz/util/sse_client"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func NewErine(model string) (*Erine, err.Error) {
	switch model {
	case domain.ModelErine4:
		return &Erine{
			ChatUrl:   erine4chatUrl,
			AppKey:    config.GetBaiduConf().AppKey,
			AppSecret: config.GetBaiduConf().AppSecret,
		}, nil
	}

	return nil, err.ModelNotSupported
}

type Erine struct {
	AppKey    string
	AppSecret string
	ChatUrl   string
}

func (*Erine) PlatformErrHandler(pErr *domain.PlatformError) string {
	if pErr.Err != nil {
		return "internal server error"
	}

	return fmt.Sprintf("%d:%s", pErr.Code, pErr.Msg)
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
	respCh := make(chan string)

	sse_client.HandleSseResp(ctx, httpResp, func(ctx context.Context, event *sse_client.SseEvent) bool {
		if event.Data == nil {
			return false
		}

		var respBody ChatCreateResp
		if err := json.Unmarshal(event.Data, &respBody); err != nil {
			hlog.CtxErrorf(ctx, "json unmarshal err: %v", err)
			errCh <- &domain.PlatformError{Err: err}
			return false
		}

		if respBody.ErrorCode != 0 || respBody.Error != "" {
			hlog.CtxErrorf(ctx, "request err: %s. %s.", respBody.ErrorMsg, respBody.ErrorDescription)
			errCh <- &domain.PlatformError{
				Code: respBody.ErrorCode,
				Msg:  respBody.ErrorMsg,
			}
			return false
		}

		if respBody.IsEnd {
			return true
		}

		if len(respBody.Result) > 0 {
			respCh <- respBody.Result
		}

		return false
	},
		func() {
			close(respCh)
			close(errCh)
		},
	)

	return respCh, errCh, nil
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
