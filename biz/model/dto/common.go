package dto

import (
	"net/http"
	"web_chat/biz/model/err"

	"github.com/cloudwego/hertz/pkg/app"
)

type CommonResp struct {
	Success bool   `json:"success"`
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func (resp *CommonResp) SetSuccess(s bool) {
	resp.Success = s
}

func (resp *CommonResp) SetCode(code int32) {
	resp.Code = code
}

func (resp *CommonResp) SetMsg(msg string) {
	resp.Message = msg
}

type RespInf interface {
	SetSuccess(s bool)
	SetCode(code int32)
	SetMsg(msg string)
}

func SuccessResp(c *app.RequestContext, resp RespInf) {
	resp.SetSuccess(true)
	resp.SetCode(err.Success.Code())
	resp.SetMsg(err.Success.Msg())

	c.JSON(http.StatusOK, resp)
}

func FailResp(c *app.RequestContext, resp RespInf, bizErr err.Error) {
	resp.SetSuccess(false)
	resp.SetCode(bizErr.Code())
	resp.SetMsg(bizErr.Msg())

	c.JSON(http.StatusOK, resp)
}
