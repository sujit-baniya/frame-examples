package main

import (
	"context"
	"fmt"
	"github.com/sujit-baniya/frame"
	"github.com/sujit-baniya/frame/pkg/protocol"
	"github.com/sujit-baniya/frame/pkg/protocol/consts"
	"net/url"
	"time"
)

func main() {
	h := server.Default(server.WithHostPorts("127.0.0.1:8080"), server.WithMaxRequestBodySize(20<<20))

	h.GET("/cookie", func(ctx context.Context, c *frame.Context) {
		mc := "myCookie"
		// get specific key
		val := c.Cookie(mc)
		if val == nil {
			// set a cookie
			fmt.Printf("There is no cookie named: %s, and make one...\n", mc)
			cookie := protocol.AcquireCookie()
			cookie.SetKey("myCookie")
			cookie.SetValue("a nice cookie!")
			cookie.SetExpire(time.Now().Add(3600 * time.Second))
			cookie.SetPath("/")
			cookie.SetHTTPOnly(true)
			cookie.SetSecure(false)
			c.Response.Header.SetCookie(cookie)
			protocol.ReleaseCookie(cookie)
			c.WriteString("A cookie is ready.")
			return
		}

		fmt.Printf("Got a cookie: %s\nAnd eat it!", val)
		// instruct upload_file to delete a cookie
		// DelClientCookie instructs the upload_file to remove the given cookie.
		// This doesn't work for a cookie with specific domain or path,
		// you should delete it manually like:
		//
		//      c := AcquireCookie()
		//      c.SetKey(mc)
		//      c.SetDomain("example.com")
		//      c.SetPath("/path")
		//      c.SetExpire(CookieExpireDelete)
		//      h.SetCookie(c)
		//      ReleaseCookie(c)
		//
		c.Response.Header.DelClientCookie(mc)

		// construct the full struct of a cookie in response's header
		respCookie := protocol.AcquireCookie()
		respCookie.SetKey(mc)
		c.Response.Header.Cookie(respCookie)
		fmt.Printf("(The expire time of cookie is set to: %s)\n", respCookie.Expire())
		protocol.ReleaseCookie(respCookie)
		c.WriteString("The cookie has been eaten.")
	})

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to url matching: /welcome?firstname=Jane&lastname=Doe&food=apple&food=fish
	h.GET("/welcome", func(ctx context.Context, c *frame.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		// shortcut for c.Request.URL.Query().Get("lastname")
		lastname := c.Query("lastname")

		// Iterate all queries and store the one with meeting the conditions in favoriteFood
		var favoriteFood []string
		c.QueryArgs().VisitAll(func(key, value []byte) {
			if string(key) == "food" {
				favoriteFood = append(favoriteFood, string(value))
			}
		})

		c.String(consts.StatusOK, "Hello %s %s, favorite food: %s", firstname, lastname, favoriteFood)
	})
	// content-type : application/x-www-form-urlencoded
	h.POST("/urlencoded", func(ctx context.Context, c *frame.Context) {
		name := c.PostForm("name")
		message := c.PostForm("message")

		c.PostArgs().VisitAll(func(key, value []byte) {
			if string(key) == "name" {
				fmt.Printf("This is %s!", string(value))
			}
		})

		c.String(consts.StatusOK, "name: %s; message: %s", name, message)
	})

	// content-type : multipart/form-data
	h.POST("/formdata", func(ctx context.Context, c *frame.Context) {
		id := c.FormValue("id")
		name := c.FormValue("name")
		message := c.FormValue("message")

		c.String(consts.StatusOK, "id: %s; name: %s; message: %s\n", id, name, message)
	})

	h.POST("/singleFile", func(ctx context.Context, c *frame.Context) {
		// single file
		file, _ := c.FormFile("file")
		fmt.Println(file.Filename)

		// Upload the file to specific dst
		c.SaveUploadedFile(file, fmt.Sprintf("./file/upload/%s", file.Filename))

		c.String(consts.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	h.POST("/multiFile", func(ctx context.Context, c *frame.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["file"]

		for _, file := range files {
			fmt.Println(file.Filename)

			// Upload the file to specific dst.
			c.SaveUploadedFile(file, fmt.Sprintf("./file/upload/%s", file.Filename))
		}
		c.String(consts.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})

	// eg. visit: http://127.0.0.1:8080/file/staticFile/main.go
	h.Static("/file", "./")

	// custom FS as you wish
	h.StaticFS("/", &frame.FS{})

	// like SimpleHTTPServer
	h.StaticFS("/try_dir", &frame.FS{Root: "./", GenerateIndexPages: true, PathRewrite: server.NewPathSlashesStripper(1)})

	h.StaticFile("/main", "./file/staticFile/main.go")

	// FileAttachment() sets the "content-disposition" header and returns the file as an "attachment".
	h.GET("/fileAttachment", func(ctx context.Context, c *frame.Context) {
		// If you use Chinese, need to encode
		fileName := url.QueryEscape("hertz")
		c.FileAttachment("./file/download/file.txt", fileName)
	})

	// File() will return the contents of the file directly
	h.GET("/file", func(ctx context.Context, c *frame.Context) {
		c.File("./file/download/file.txt")
	})
	h.Spin()
}
