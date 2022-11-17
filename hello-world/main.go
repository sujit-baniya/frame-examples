package main

import (
	"context"
	"github.com/sujit-baniya/frame"
	"github.com/sujit-baniya/frame/pkg/common/utils"
	"github.com/sujit-baniya/frame/pkg/protocol/consts"
	"github.com/sujit-baniya/frame/server"
)

type Test struct {
	A string
	B string
}

func main() {
	h := server.Default()

	h.GET("/ping", func(c context.Context, ctx *frame.Context) {
		ctx.JSON(consts.StatusOK, utils.H{"ping": "pong"})
	})

	h.Spin()
}
