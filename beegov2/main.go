package main

import (
	_ "bee2/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

