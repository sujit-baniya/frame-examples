package main

import (
	"context"
	"fmt"
	"github.com/sujit-baniya/frame/client"
	"github.com/sujit-baniya/frame/pkg/protocol"
	"github.com/sujit-baniya/frame/pkg/protocol/consts"
	"github.com/sujit-baniya/frame/pkg/server/binding"
	"time"
)

func init() {
	binding.SetLooseZeroMode(false)
}

func main() {
	h := server.Default(server.WithHostPorts("127.0.0.1:8080"))

	h.GET("looseZero", func(ctx context.Context, c *frame.Context) {
		type Loose struct {
			A int `query:"a"`
		}
		var req Loose
		err := c.BindAndValidate(&req)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			panic(err)
		}
		fmt.Printf("req: %v\n", req)
	})

	go h.Spin()

	time.Sleep(1000 * time.Millisecond)
	c, _ := client.NewClient()
	req := protocol.Request{}
	resp := protocol.Response{}
	req.SetMethod(consts.MethodGet)
	req.SetRequestURI("http://127.0.0.1:8080/looseZero?a")
	c.Do(context.Background(), &req, &resp)
}
