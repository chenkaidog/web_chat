package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func LoginPage(ctx context.Context, c *app.RequestContext) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func HomePage(ctx context.Context, c *app.RequestContext) {
	c.HTML(http.StatusOK, "home.html", nil)
}

func ReadMePage(ctx context.Context, c *app.RequestContext) {
	c.HTML(http.StatusOK, "readme.html", utils.H{
		"readme": readme,
	})
}

func ChatPage(ctx context.Context, c *app.RequestContext) {
	chatType, _ := c.GetQuery("chat_type") 

	options := []struct {
		Value string
		Label string
	}{
		{"ernie4", "文心4"},
		{"gpt3", "gpt-3.5"},
		{"gpt4", "gpt-4"},
	}

	c.HTML(http.StatusOK, "chat.html", utils.H{
		"Options": options,
		"ChatType": chatType,
	})
}

func PasswordUpdatePage(ctx context.Context, c *app.RequestContext)  {
	c.HTML(http.StatusOK, "update_password.html", nil)
}

var readme = `
# 使用声明

## 账号相关
1. 每个账号只能在最多3台设备上登陆，最早登陆的账号会被挤占下线。
2. 为防止网站被攻击，网站不提供注册功能。本站点不对外开放，各位用户请保管好自己的账号，不要泄漏密码。
3. 请不要进行频繁的请求，否则会被安全策略自动封锁IP或者账号。

## AI对话原理
网站不会存储用户对话的记录，对话记录存储在用户的浏览器本地，各大AI对话平台也声称不会存储用户的对话。用户的提问和AI的响应会构成上下文，AI正是根据这个上下文进行“思考”的。既然平台不会存储用户的对话，那么这个上下文放在哪里呢？答案是用户每次进行AI对话的时候都会在请求中带上之前的聊天记录，这个聊天记录就是上下文，而这个上下文只会保存在用户的浏览器中。

AI都是大数据训练出来的，它们不会了解最新的知识。

AI每次响应都会消耗一定的tokens，tokens可以理解为AI分析的字数，包括请求中的字数和响应中的字数。因此随着上下文不断增长，请求消耗的token也会逐步增长。每个平台的tokens限额都是要充值的。

## 对话类型
对话分为了无记忆对话和有记忆对话。

无记忆对话中不会带上用户的上下文，每次响应都只是针对用户的单次提问，用户不能继续追问。这样只是为了节省成本，随着上下文的增长，每次请求都会消耗更多的token。

有记忆对话就是各个平台所展示的对话功能，用户在AI响应后能继续追问，从而达到更好的响应效果。

## AI模型
目前本站点支持百度的文心4和openai的gpt3.5以及gpt4。
这几个模型是在中文语境下评分最高的模型: gpt3.5<文心4<gpt4
`
