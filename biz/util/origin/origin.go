package origin

import (
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/mssola/user_agent"
)

func GetIp(c *app.RequestContext) string {
	host := string(c.RemoteAddr().String())
	idx := strings.Index(host, ":")
	if idx > 0 {
		return host[:idx]
	}

	return host
}

func GetDevice(c *app.RequestContext) string {
	userAgent := user_agent.New(string(c.UserAgent()))
	if name, _ := userAgent.Browser(); name != "" {
		return name
	}

	return "UNKNOWN"
}