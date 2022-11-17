package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// Maximum wait time before exit, if not specified the default is 5s
	h := server.Default(server.WithHostPorts("127.0.0.1:8080"), server.WithExitWaitTime(3*time.Second))

	h.OnShutdown = append(h.OnShutdown, func(ctx context.Context) {
		fmt.Println("hook 1")
		<-ctx.Done()
		fmt.Println("exit timeout!")
	})

	h.OnShutdown = append(h.OnShutdown, func(ctx context.Context) {
		fmt.Println("hook 2")
	})

	h.Spin()
}
