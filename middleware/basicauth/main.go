package main

import (
	"context"
	"github.com/sujit-baniya/frame/pkg/app"
	"github.com/sujit-baniya/frame/pkg/app/middlewares/server/basic_auth"
	"github.com/sujit-baniya/frame/pkg/app/server"
	"github.com/sujit-baniya/frame/pkg/protocol/consts"
)

func main() {
	h := server.Default(server.WithHostPorts("127.0.0.1:8080"))
	h.Use(basic_auth.BasicAuth(map[string]string{
		"test1": "value1",
		"test2": "value2",
	}))

	h.GET("/basicAuth", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "hello hertz")
	})

	h.Spin()
}
