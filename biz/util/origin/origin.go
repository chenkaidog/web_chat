package origin

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/mssola/user_agent"
)

func GetIp(c *app.RequestContext) string {
	return c.Request.Header.Get("X-Real-IP") // 从nginx透传的IP
}

func GetDevice(c *app.RequestContext) string {
	userAgent := user_agent.New(string(c.UserAgent()))
	if name, _ := userAgent.Browser(); name != "" {
		return name
	}

	return "UNKNOWN"
}
