package main

import (
	"context"
	"fmt"
	"github.com/sujit-baniya/frame/pkg/app"
	"github.com/sujit-baniya/frame/pkg/app/client"
	"github.com/sujit-baniya/frame/pkg/app/server"
	"github.com/sujit-baniya/frame/pkg/app/server/binding"
	"github.com/sujit-baniya/frame/pkg/protocol"
	"github.com/sujit-baniya/frame/pkg/protocol/consts"
	"time"
)

type BindError struct {
	ErrType, FailField, Msg string
}

// Error implements error interface.
func (e *BindError) Error() string {
	if e.Msg != "" {
		return e.ErrType + ": expr_path=" + e.FailField + ", cause=" + e.Msg
	}
	return e.ErrType + ": expr_path=" + e.FailField + ", cause=invalid"
}

type ValidateError struct {
	ErrType, FailField, Msg string
}

// Error implements error interface.
func (e *ValidateError) Error() string {
	if e.Msg != "" {
		return e.ErrType + ": expr_path=" + e.FailField + ", cause=" + e.Msg
	}
	return e.ErrType + ": expr_path=" + e.FailField + ", cause=invalid"
}

func init() {
	CustomBindErrFunc := func(failField, msg string) error {
		err := BindError{
			ErrType:   "bindErr",
			FailField: "[bindFailField]: " + failField,
			Msg:       "[bindErrMsg]: " + msg,
		}

		return &err
	}

	CustomValidateErrFunc := func(failField, msg string) error {
		err := ValidateError{
			ErrType:   "validateErr",
			FailField: "[validateFailField]: " + failField,
			Msg:       "[validateErrMsg]: " + msg,
		}

		return &err
	}

	binding.SetErrorFactory(CustomBindErrFunc, CustomValidateErrFunc)
}

func main() {
	h := server.Default(server.WithHostPorts("127.0.0.1:8080"))

	h.GET("bindErr", func(ctx context.Context, c *app.RequestContext) {
		type TestBind struct {
			A string `query:"a,required"`
		}
		var req TestBind
		err := c.Bind(&req)
		fmt.Printf("error: %v\n", err)
	})

	h.GET("validateErr", func(ctx context.Context, c *app.RequestContext) {
		type TestValidate struct {
			B int `query:"b" vd:"$>100; msg:'C must greater than 100'"`
		}
		var req TestValidate
		err := c.Bind(&req)
		if err != nil {
			panic(err)
		}
		err = c.Validate(&req)
		fmt.Printf("error: %v\n", err)
	})

	go h.Spin()

	time.Sleep(1000 * time.Millisecond)
	c, _ := client.NewClient()
	req := protocol.Request{}
	resp := protocol.Response{}
	req.SetMethod(consts.MethodGet)
	req.SetRequestURI("http://127.0.0.1:8080/bindErr")
	c.Do(context.Background(), &req, &resp)

	req.SetRequestURI("http://127.0.0.1:8080/validateErr?b=1")
	c.Do(context.Background(), &req, &resp)
}
