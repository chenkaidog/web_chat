package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"web_chat/biz/config"
	"web_chat/biz/model/domain"
	"web_chat/biz/model/err"
	"web_chat/biz/util/sse_client"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func NewChatGPT(model string) (*ChatGPT, err.Error) {
	switch model {
	case domain.ModelGPT3:
		return &ChatGPT{
			Model:  modelGPT3,
			ApiKey: config.GetOpenAIConf().ApiKey,
		}, nil
	case domain.ModelGPT4:
		return &ChatGPT{
			Model:  modelGPT4,
			ApiKey: config.GetOpenAIConf().ApiKey,
		}, nil
	}

	return nil, err.ModelNotSupported
}

type ChatGPT struct {
	Model  string
	ApiKey string
}

func (*ChatGPT) PlatformErrHandler(pErr *domain.PlatformError) string {
	if pErr.Err != nil {
		return "internal server error"
	}

	return fmt.Sprintf("%d:%s", pErr.Code, pErr.Msg)
}

func (gpt *ChatGPT) StreamChat(ctx context.Context, chatContext []*domain.ChatContent) (chan string, chan *domain.PlatformError, error) {
	req, err := gpt.newStreamChatRequest(ctx, chatContext)
	if err != nil {
		return nil, nil, err
	}

	httpClient, err := newHttpClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		hlog.CtxErrorf(ctx, "http request err: %v", err)
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		respContent, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return nil, nil, err
		}

		hlog.CtxErrorf(ctx, "status_code: %d, error_msg: %s", resp.StatusCode, respContent)
		return nil, nil, errors.New("request fails")
	}

	errCh := make(chan *domain.PlatformError)
	respCh := make(chan string)

	sse_client.HandleSseResp(ctx, resp, func(ctx context.Context, event *sse_client.SseEvent) bool {
		if event.Data == nil {
			return false
		}

		var respBody ChatCreateResp
		if err := json.Unmarshal(event.Data, &respBody); err != nil {
			hlog.CtxErrorf(ctx, "json unmarshal err: %v", err)
			errCh <- &domain.PlatformError{Err: err}
			return false
		}

		// todo: handle biz err

		if len(respBody.Choices) > 0 {
			choice := respBody.Choices[0]
			if choice.FinishReason == finishReasonStop {
				return true
			}

			if len(choice.Delta.Content) > 0 {
				respCh <- choice.Delta.Content
			}
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

func (gpt *ChatGPT) newStreamChatRequest(ctx context.Context, chatContext []*domain.ChatContent) (*http.Request, error) {
	var messages []Message
	for _, content := range chatContext {
		messages = append(
			messages,
			Message{
				Role:    roleMapper[content.Role],
				Content: content.Content,
			},
		)
	}

	reqBody, err := json.Marshal(
		&ChatCreateReq{
			Model:    gpt.Model,
			Stream:   true,
			Messages: messages,
		},
	)
	if err != nil {
		hlog.CtxErrorf(ctx, "json marshal err: %v", err)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, chatUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		hlog.CtxErrorf(ctx, "new http request err: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+gpt.ApiKey)

	return req, nil
}

func newHttpClient(ctx context.Context) (*http.Client, error) {
	proxyUrl, err := url.Parse(proxyUrl)
	if err != nil {
		hlog.CtxErrorf(ctx, "parse proxy err: %v", err)
		return nil, err
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	return httpClient, nil
}
