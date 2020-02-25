package main

import (
	"github.com/xiui/zyh"
)

func main() {

	r := zyh.Default()

	r.GET("/test", func(ctx *zyh.Context) {

		ctx.JSON(200, map[string]string{
			"msg":"ok",
		})
	})

	r.Use(func(ctx *zyh.Context) {
		ctx.String(200, "r.middle ")
		ctx.Next()
	})

	g := r.Group("/v1")

	g.Use(func(ctx *zyh.Context) {
		ctx.String(200, "g.middle ")
		ctx.Next()
	})

	g.GET("/test", func(ctx *zyh.Context) {
		ctx.String(200, "g.get")
	})

	g2 := r.Group("/v2")

	g2.Use(func(ctx *zyh.Context) {
		ctx.String(200, "g2.middle ")
		ctx.Next()
	})

	g2.GET("/test", func(ctx *zyh.Context) {
		ctx.String(200, "g2.get")
	})



	r.Run(":8080")
}
