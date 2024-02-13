package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
)

func StreamChat(ctx context.Context, c *app.RequestContext) {
	c.SetStatusCode(http.StatusOK)
	// s := sse.NewStream(c)

}
