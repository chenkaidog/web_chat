package dto

type Message struct {
	Role    Role   `form:"role" json:"role" binding:"required"`
	Content string `form:"content" json:"content" binding:"required,max=500"`
}

type ChatCreateReq struct {
	Platform Platform  `json:"platform" binding:"required"`
	Model    Model     `json:"model" binding:"required"`
	Messages []Message `json:"messages" binding:"required,max=20"`
}

type ChatCreateResp struct {
	CommonResp

	CreatedAt int64  `json:"created_at,omitempty"`
	IsEnd     bool   `json:"is_end,omitempty"`
	Content   string `json:"content,omitempty"`
}

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Platform string

const (
	PlatformBaidu  Platform = "baidu"
	PlatformOpenAI Platform = "openai"
)

type Model string

const (
	ModelErine4 = "erine-4"
	ModelGPT3   = "gpt-3.5"
	ModelGPT4   = "gpt-4"
)
