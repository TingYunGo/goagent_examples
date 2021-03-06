package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	tingyun "github.com/TingYunGo/goagent"
)

func DoWorkerJob() error {
	//听云API嵌码 -> 取当前协程上的事务追踪对象,
	//               创建一个组件追踪对象,追踪本函数的执行性能和错误
	component := tingyun.GetAction().CreateComponent("DoWorkerJob")
	//听云API嵌码 -> 在函数结束时,调用tingyun.Component.Finish方法 结束组件追踪
	defer component.Finish()
	time.Sleep(time.Millisecond * 100)
	r := errors.New("test error")
	//听云API嵌码 -> 组件执行过程中有错误出现时,
	//               调用tingyun.Component.SetError方法通知组件追踪对象记录错误
	component.SetError(r, "test", 1)
	return r
}
func main() {
	listen := ":3000"
	if len(os.Args) > 1 {
		listen = os.Args[1]
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().Format("2006-01-02 15:04:05.000")
		fmt.Println(now, "RID:", tingyun.GetGID(), "URI:", r.URL.String())
		header := w.Header()
		header.Set("Cache-Control", "no-cache")
		header.Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		DoWorkerJob()
		b, _ := json.Marshal(map[string]interface{}{
			"status": "success",
			"URI":    r.URL.String(),
			"GID":    tingyun.GetGID(),
		})
		w.Write(b)
	})
	fmt.Println("Service Listen @", listen)
	http.ListenAndServe(listen, nil)
}
