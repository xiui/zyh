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

	r.Run(":8080")
}
