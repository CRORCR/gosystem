package response

import (
	"github.com/kataras/iris"

	"gosystem/bootstrap"
	"gosystem/conf"
)

func Configure(b *bootstrap.Bootstrapper) {
	b.Use(func(ctx iris.Context) {
		ctx.Values().Set("result", &conf.Result{
			Code: 0,
			Msg:  "",
			Data: nil,
		})
		ctx.Next()
	})

	b.Done(func(ctx iris.Context) {
		ctx.JSON(ctx.Values().Get("result"))
	})
}
