package routers

import (
	"bee2/controllers"
	"encoding/json"
	"net/http"


	//go get github.com/astaxie/beego
	"github.com/beego/beego/v2/server/web/context"

	beego "github.com/beego/beego/v2/server/web"
)

func handler(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("Cache-Control", "no-cache")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(map[string]interface{}{
		"status": "success",
		"result": r.Host,
	})
	//	fmt.Printf("on handler\n")
	w.Write(b)
}

type MyHandlerWrapper func(http.ResponseWriter, *http.Request)

// ServeHTTP calls f(w, r).
func (f MyHandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

type MyHandler struct {
}

// ServeHTTP calls f(w, r).
func (f *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler(w, r)
}

func init() {
	beego.Any("/api/:id([0-9]+)", func(ctx *context.Context) {
		ctx.Output.Body([]byte("bar"))
	})
	beego.Handler("/handler", http.HandlerFunc(handler))
	beego.Handler("/handler1", MyHandlerWrapper(handler))
	beego.Handler("/handler2", &MyHandler{})
	beego.Router("/", &controllers.MainController{})
}
