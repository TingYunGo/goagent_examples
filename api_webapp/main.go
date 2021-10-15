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
	//手动嵌码
	component := tingyun.GetAction().CreateComponent("DoWorkerJob")
	defer component.Finish()

	time.Sleep(time.Millisecond * 100)
	r := errors.New("test error")
	component.SetError(r, "test", 1)
	return r
}
func main() {
	listen := ":3000"
	if len(os.Args) > 1 {
		listen = os.Args[1]
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"), "RID:", tingyun.GetGID(), "URI:", r.URL.String())
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
