package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	tingyun "github.com/TingYunGo/goagent"
)

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
