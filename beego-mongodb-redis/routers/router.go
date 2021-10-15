package routers

import (
	"beego-mongodb-redis/controllers"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	beego.Router("/*", &controllers.MainController{})
	beego.Get("/favicon.ico", func(ctx *context.Context) {
		http.ServeFile(ctx.ResponseWriter, ctx.Request, "static/img/favicon.ico")
	})
}
