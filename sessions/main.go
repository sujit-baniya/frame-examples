package main

import (
	"context"
	"github.com/sujit-baniya/frame"
	"github.com/sujit-baniya/frame/middlewares/server/sessions"
	"github.com/sujit-baniya/frame/middlewares/server/sessions/cookie"
	"github.com/sujit-baniya/frame/middlewares/server/sessions/redis"
	"github.com/sujit-baniya/frame/pkg/common/utils"
	"github.com/sujit-baniya/frame/server"
)

func main() {
	r()
}

func c() {
	h := server.New(server.WithHostPorts(":8080"))
	store := cookie.NewStore([]byte("secret"))
	h.Use(sessions.Sessions("mysession", store))
	h.GET("/incr", func(ctx context.Context, c *frame.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.JSON(200, utils.H{"count": count})
	})
	h.Spin()
}

func r() {
	h := server.Default(server.WithHostPorts(":8000"))
	store, _ := redis.NewStore("localhost:6379", "", 0)
	h.Use(sessions.Sessions("mysession", store))

	h.GET("/incr", func(ctx context.Context, c *frame.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.JSON(200, utils.H{"count": count})
	})
	h.Spin()
}
