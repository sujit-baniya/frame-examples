package main

import (
	"context"
	"github.com/sujit-baniya/frame"
	"github.com/sujit-baniya/frame/middlewares/server/basic_auth"
	"github.com/sujit-baniya/frame/pkg/protocol/consts"
	"github.com/sujit-baniya/frame/server"
)

func main() {
	h := server.Default(server.WithHostPorts("127.0.0.1:8080"))

	h.GET("/", func(ctx context.Context, c *frame.Context) {
		c.String(consts.StatusOK, "hello hertz")
	})
	h.Use(basic_auth.BasicAuth(map[string]string{
		"test1": "value1",
		"test2": "value2",
	}))

	h.GET("/basicAuth", func(ctx context.Context, c *frame.Context) {
		c.String(consts.StatusOK, "hello hertz")
	})

	h.Spin()
}
