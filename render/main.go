package main

import (
	"context"
	"fmt"
	"github.com/sujit-baniya/frame/pkg/app"
	"github.com/sujit-baniya/frame/pkg/app/server"
	"github.com/sujit-baniya/frame/pkg/common/utils"
	"github.com/sujit-baniya/frame/pkg/protocol/consts"
	"net/http"
	"time"
)

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}

func main() {
	h := server.Default(server.WithHostPorts(":8080"))
	h.SetHTMLTemplate("./resources", ".html")

	h.GET("/index", func(c context.Context, ctx *app.RequestContext) {
		ctx.HTML(http.StatusOK, "index", utils.H{
			"title": "Main website",
		})
	})

	// utils.H is a shortcut for map[string]interface{}
	h.GET("/someJSON", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"message": "hey", "status": consts.StatusOK})
	})

	h.GET("/moreJSON", func(ctx context.Context, c *app.RequestContext) {
		// You also can use a struct
		var msg struct {
			Company  string `json:"company"`
			Location string
			Number   int
		}
		msg.Company = "company"
		msg.Location = "location"
		msg.Number = 123
		// Note that msg.Company becomes "company" in the JSON
		// Will output  :   {"company": "company", "Location": "location", "Number": 123}
		c.JSON(consts.StatusOK, msg)
	})

	h.GET("/pureJson", func(ctx context.Context, c *app.RequestContext) {
		c.PureJSON(consts.StatusOK, utils.H{
			"html": "<p> Hello World </p>",
		})
	})

	h.GET("/someData", func(ctx context.Context, c *app.RequestContext) {
		c.Data(consts.StatusOK, "text/plain; charset=utf-8", []byte("hello"))
	})

	h.GET("/externalRedirect", func(ctx context.Context, c *app.RequestContext) {
		c.Redirect(consts.StatusMovedPermanently, []byte("http://www.google.com/"))
	})

	h.GET("/internalRedirect", func(ctx context.Context, c *app.RequestContext) {
		c.Redirect(consts.StatusFound, []byte("/foo"))
	})

	h.GET("/foo", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "hello, world")
	})

	h.Spin()
}
