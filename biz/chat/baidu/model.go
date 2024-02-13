package baidu

import "web_chat/biz/model/domain"

const (
	accessUrl     = "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials"
	erine4chatUrl = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions_pro"
)

const (
	roleAssistant = "assistant"
	roleUser      = "user"
)

var roleMapper = map[domain.Role]string{
	domain.RoleUser:      roleUser,
	domain.RoleAssistant: roleAssistant,
}

type AppAccessInfo struct {
	AccessToken   string `json:"access_token"`
	ExpiresIn     int64  `json:"expires_in"`
	ErrorCode     string `json:"error"`
	Description   string `json:"error_description"`
	SessionKey    string `json:"session_key"`
	RefreshToken  string `json:"refresh_token"`
	Scope         string `json:"scope"`
	SessionSecret string `json:"session_secret"`
}

// https://cloud.baidu.com/doc/WENXINWORKSHOP/s/clntwmv7t#header%E5%8F%82%E6%95%B0
type ChatCreateReq struct {
	Messages       []*Message  `json:"messages"`
	Functions      []*Function `json:"functions,omitempty"`
	Temperature    float64     `json:"temperature,omitempty"`
	TopP           float64     `json:"top_p,omitempty"`
	PenaltyScore   float64     `json:"penalty_score,omitempty"`
	Stream         bool        `json:"stream,omitempty"`
	System         string      `json:"system,omitempty"`
	Stop           []string    `json:"stop,omitempty"`
	DisableSearch  bool        `json:"disable_search,omitempty"`
	EnableCitation bool        `json:"enable_citation,omitempty"`
	ResponseFormat string      `json:"response_foramt,omitempty"`
	UserId         string      `json:"user_id,omitempty"`
}

type ChatCreateResp struct {
	// access error
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`

	// bizError
	ErrorDescription string `json:"error_description"`
	Error            string `json:"error"`

	ID               string        `json:"id"`
	Object           string        `json:"object"`
	Created          int           `json:"created"`
	SentenceID       int           `json:"sentence_id"`
	IsEnd            bool          `json:"is_end"`
	IsTruncated      bool          `json:"is_truncated"`
	FinishReason     string        `json:"finish_reason"`
	SearchInfo       *SearchInfo   `json:"search_info"`
	Result           string        `json:"result"`
	NeedClearHistory bool          `json:"need_clear_history"`
	Flag             int           `json:"flag"`
	BanRound         int           `json:"ban_round"`
	Usage            *Usage        `json:"usage"`
	FunctionCall     *FunctionCall `json:"function_call"`
}

type Message struct {
	Role         string        `json:"role"`
	Content      string        `json:"content"`
	Name         string        `json:"name,omitempty"`
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
}

type Function struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"`
	Responses   interface{} `json:"responses,omitempty"`
	Examples    []*Example  `json:"examples,omitempty"`
}

type Example struct {
	Role         string        `json:"role"`
	Content      string        `json:"content"`
	Name         string        `json:"name,omitempty"`
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
}

type FunctionCall struct {
	Name     string `json:"name"`
	Args     string `json:"arguments"`
	Thoughts string `json:"thoughts,omitempty"`
}

type SearchInfo struct {
	SearchResults []*SearchResult `json:"search_results"`
}

type SearchResult struct {
	Index int    `json:"index"`
	Url   string `json:"url"`
	Title string `json:"title"`
}

type Usage struct {
	PromptToken      int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
