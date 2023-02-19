package main

import (
	"day4/gee"
	"net/http"
)

func main() {
	engine := gee.NewEngine()
	engine.GET("/index", func(c *gee.Context) {
		c.String(http.StatusOK, "<h1>Index page</h1>")
	})
	v1 := engine.Group("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee </h1>")
		})
		v1.GET("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := engine.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{"username": c.PostForm("username"),
				"password": c.PostForm("password")})
		})
	}
	engine.Run(":9999")
}
