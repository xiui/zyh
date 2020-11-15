package main

import (
	"fmt"
	"github.com/xiui/zyh"
	"time"
)

func main() {

	r := zyh.Default()



	r.GET("/test", func(ctx *zyh.Context) {

		ctx.Redirect("https://www.baidu.com", 302)
	})

	r.POST("/test", func(ctx *zyh.Context) {

		fmt.Println(ctx.Params)


	})

	r.UseMiddleware(func(ctx *zyh.Context) {
		ctx.String(200, "r.middle1 ")
		ctx.Next()
	}, func(ctx *zyh.Context) {
		ctx.String(200, "r.middle2 ")
		ctx.Next()
	})

	r.AddMiddleware(func(ctx *zyh.Context) {
		ctx.String(200, "r.middle add 1 ")

		fmt.Println(time.Now().Unix())
		ctx.Next()
		fmt.Println(time.Now().Unix())
	})

	g := r.Group("/v1")

	g.UseMiddleware(func(ctx *zyh.Context) {
		ctx.String(200, "g.middle ")
		time.Sleep(5 * time.Second)
		ctx.Next()
	})

	g.GET("/test", func(ctx *zyh.Context) {
		ctx.String(200, "g.get")
	})

	g2 := r.Group("/v2")

	g2.UseMiddleware(func(ctx *zyh.Context) {
		ctx.String(200, "g2.middle1 \n")
		ctx.Next()
	}, func(ctx *zyh.Context) {
		ctx.String(200, "g2.middle2 \n")
		ctx.Next()
	})

	g2.GET("/test", func(ctx *zyh.Context) {
		ctx.String(200, "g2.get1\n")

		fmt.Println("g2.get1")
		ctx.Next()

	}, func(ctx *zyh.Context) {
		ctx.String(200, "g2.get2")
		fmt.Println("g2.get2")

	})



	r.Run(":8080")
}
