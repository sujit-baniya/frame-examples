package main

import (
	"context"
	"fmt"
	"github.com/sujit-baniya/frame"
	"github.com/sujit-baniya/frame/client"
	"github.com/sujit-baniya/frame/pkg/protocol"
	"github.com/sujit-baniya/frame/pkg/protocol/consts"
	"github.com/sujit-baniya/frame/server/binding"
	"time"
)

type ValidateStruct struct {
	A string `query:"a" vd:"test($)"`
}

func init() {
	binding.MustRegValidateFunc("test", func(args ...interface{}) error {
		if len(args) != 1 {
			return fmt.Errorf("the args must be one")
		}
		s, _ := args[0].(string)
		if s == "123" {
			return fmt.Errorf("the args can not be 123")
		}
		return nil
	})
}

func main() {
	h := server.Default(server.WithHostPorts("127.0.0.1:8080"))

	h.GET("customValidate", func(ctx context.Context, c *frame.Context) {
		var req ValidateStruct
		err := c.Bind(&req)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		err = c.Validate(&req)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	})

	go h.Spin()

	time.Sleep(1000 * time.Millisecond)
	c, _ := client.NewClient()
	req := protocol.Request{}
	resp := protocol.Response{}
	req.SetMethod(consts.MethodGet)
	req.SetRequestURI("http://127.0.0.1:8080/customValidate?a=123")
	c.Do(context.Background(), &req, &resp)
}
