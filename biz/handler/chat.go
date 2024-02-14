package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"web_chat/biz/chat"
	"web_chat/biz/model/domain"
	"web_chat/biz/model/dto"
	"web_chat/biz/model/err"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/sse"
)

var roleMapper = map[dto.Role]domain.Role{
	dto.RoleSystem:    domain.RoleSystem,
	dto.RoleUser:      domain.RoleUser,
	dto.RoleAssistant: domain.RoleAssistant,
}

var platformMapper = map[dto.Platform]string{
	dto.PlatformBaidu:  domain.PlatformBaidu,
	dto.PlatformOpenAI: domain.PlatformOpenai,
}

var modelMapper = map[dto.Model]string{
	dto.ModelErine4: domain.ModelErine4,
	dto.ModelGPT3:   domain.ModelGPT3,
	dto.ModelGPT4:   domain.ModelGPT4,
}

func StreamChat(ctx context.Context, c *app.RequestContext) {
	var stdErr error
	var req dto.ChatCreateReq
	var resp dto.ChatCreateResp
	if stdErr = c.BindAndValidate(&req); stdErr != nil {
		hlog.CtxErrorf(ctx, "BindAndValidate fail, %v", stdErr)
		dto.FailResp(c, &resp, err.ParamError)
		return
	}

	chatImpl, bizErr := chat.NewchatImpl(
		platformMapper[req.Platform],
		modelMapper[req.Model],
	)
	if bizErr != nil {
		hlog.CtxErrorf(ctx, "request param invalid[%v]: %s=>%s", bizErr, req.Platform, req.Model)
		dto.FailResp(c, &resp, err.ParamError)
		return
	}

	var chatContext []*domain.ChatContent
	for _, msg := range req.Messages {
		chatContext = append(
			chatContext,
			&domain.ChatContent{
				Role:    roleMapper[msg.Role],
				Content: msg.Content,
			},
		)
	}

	cancelCtx, cancelFunx := context.WithCancel(ctx)
	defer cancelFunx()

	respCh, errCh, stdErr := chatImpl.StreamChat(cancelCtx, chatContext)
	if stdErr != nil {
		hlog.CtxErrorf(cancelCtx, "chat fail: %v", stdErr)
		dto.FailResp(c, &resp, err.InternalServerError)
		return
	}

	c.SetStatusCode(http.StatusOK)
	ssePublisher := sse.NewStream(c)
	timeout := time.NewTimer(time.Second * 30)

	for {
		select {
		case <-timeout.C:
			timeoutPublish(ssePublisher)
			return
		case msg, ok := <-respCh:
			if !ok {
				endPublish(ssePublisher)
				return
			}
			timeout.Reset(time.Second * 30)
			contentPublish(ssePublisher, msg)
		case pErr, ok := <-errCh:
			if ok {
				errorPublish(ssePublisher, chatImpl.PlatformErrHandler(pErr))
				return
			}
		}
	}
}

func timeoutPublish(stream *sse.Stream) {
	resp := &dto.ChatCreateResp{
		CommonResp: dto.CommonResp{
			Success: false,
			Code:    err.ResponseTimeoutError.Code(),
			Message: err.ResponseTimeoutError.Msg(),
		},
		IsEnd:     true,
		CreatedAt: time.Now().Unix(),
	}

	data, _ := json.Marshal(resp)

	stream.Publish(
		&sse.Event{
			Data: data,
		},
	)
}

func endPublish(stream *sse.Stream) {
	resp := &dto.ChatCreateResp{
		CommonResp: dto.CommonResp{
			Success: true,
		},
		IsEnd:     true,
		CreatedAt: time.Now().Unix(),
	}

	data, _ := json.Marshal(resp)

	stream.Publish(
		&sse.Event{
			Data: data,
		},
	)
}

func contentPublish(stream *sse.Stream, content string) {
	resp := &dto.ChatCreateResp{
		CommonResp: dto.CommonResp{
			Success: true,
		},
		Content:   content,
		CreatedAt: time.Now().Unix(),
	}

	data, _ := json.Marshal(resp)

	stream.Publish(
		&sse.Event{
			Data: data,
		},
	)
}

func errorPublish(stream *sse.Stream, content string) {
	resp := &dto.ChatCreateResp{
		CommonResp: dto.CommonResp{
			Success: false,
			Code:    err.PlatformError.Code(),
			Message: content,
		},
		IsEnd:     true,
		Content:   content,
		CreatedAt: time.Now().Unix(),
	}

	data, _ := json.Marshal(resp)

	stream.Publish(
		&sse.Event{
			Data: data,
		},
	)
}
