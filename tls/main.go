package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/sujit-baniya/frame"
	"github.com/sujit-baniya/frame/server"
	"time"

	"github.com/sujit-baniya/frame/client"
	"github.com/sujit-baniya/frame/pkg/network/standard"
	"github.com/sujit-baniya/frame/pkg/protocol"

	"github.com/sujit-baniya/frame/pkg/protocol/consts"
)

func main() {
	cfg := &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}
	cert, err := tls.LoadX509KeyPair("./frame.crt", "./frame.key")
	if err != nil {
		fmt.Println(err.Error())
	}
	cfg.Certificates = append(cfg.Certificates, cert)

	h := server.Default(server.WithTLS(cfg), server.WithHostPorts("127.0.0.1:8443"))

	h.Use(func(c context.Context, ctx *frame.Context) {
		fmt.Fprint(ctx, "Before real handle...\n")
		ctx.Next(c)
		fmt.Fprint(ctx, "After real handle...\n")
	})

	h.GET("/ping", func(c context.Context, ctx *frame.Context) {
		ctx.String(consts.StatusOK, "TLS test\n")
	})

	h.Spin()
}

func doTlsRequest() {
	clientCfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	c, err := client.NewClient(
		client.WithTLSConfig(clientCfg),
		client.WithDialer(standard.NewDialer()),
	)
	if err != nil {
		fmt.Println(err.Error())
	}
	req, res := protocol.AcquireRequest(), protocol.AcquireResponse()
	defer func() {
		protocol.ReleaseRequest(req)
		protocol.ReleaseResponse(res)
	}()
	req.SetMethod(consts.MethodGet)                            // set request method
	req.Header.SetContentTypeBytes([]byte("application/json")) // set request header
	req.SetRequestURI("https://localhost:8443/ping")           // set request url
	err = c.Do(context.Background(), req, res)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%v\n", string(res.Body())) // read response body
	time.Sleep(time.Second)
}
