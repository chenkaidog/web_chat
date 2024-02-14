package chat

import (
	"context"
	"web_chat/biz/chat/baidu"
	"web_chat/biz/chat/openai"
	"web_chat/biz/model/domain"
	"web_chat/biz/model/err"
)

type ChatInf interface {
	StreamChat(ctx context.Context, chatContext []*domain.ChatContent) (chan string, chan *domain.PlatformError, error)
	PlatformErrHandler(pErr *domain.PlatformError) string
}

func NewchatImpl(platform, model string) (ChatInf, err.Error) {
	switch platform {
	case domain.PlatformBaidu:
		return baidu.NewErine(model)
	case domain.PlatformOpenai:
		return openai.NewChatGPT(model)
	}

	return nil, err.PlatformNotSupported
}
