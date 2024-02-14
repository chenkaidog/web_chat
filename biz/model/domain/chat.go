package domain

type ChatContent struct {
	Role    Role
	Content string
}

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type PlatformError struct {
	Err  error
	Code int
	Msg  string
}

const (
	PlatformBaidu  = "baidu"
	ModelErine4    = "erine-4"
	PlatformOpenai = "openai"
	ModelGPT4      = "gpt-4"
	ModelGPT3      = "gpt-3.5"
)
