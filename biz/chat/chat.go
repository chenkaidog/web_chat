package chat

import (
	"context"
	"web_chat/biz/chat/baidu"
	"web_chat/biz/chat/openai"
	"web_chat/biz/model/domain"
)

type ChatInf interface {
	StreamChat(ctx context.Context, chatContext []*domain.ChatContent) (chan string, chan *domain.PlatformError, error)
}

func NewchatImpl(platform, model string) ChatInf {
	switch platform {
	case domain.PlatformBaidu:
		return baidu.NewErine(model)
	case domain.PlatformOpenai:
		return openai.NewChatGPT(model)
	}

	return nil
}
