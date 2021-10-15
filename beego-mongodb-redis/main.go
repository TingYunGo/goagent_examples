package main

import (
	"beego-mongodb-redis/controllers"
	_ "beego-mongodb-redis/routers"
	"fmt"
	"os"

	"github.com/TingYunGo/libgo/ascii"

	_ "github.com/beego/beego/v2/server/web/session/redis"

	beego "github.com/beego/beego/v2/server/web"
)

var redisunc = ""

func getRedisInfo() (host, passwd string) {
	offset := ascii.StrRChr(redisunc, '@')
	if offset == -1 {
		return redisunc, ""
	}
	return redisunc[offset+1:], ascii.SubString(redisunc, 0, offset)
}
func main() {

	redisunc = os.Getenv("REDIS")
	if len(redisunc) == 0 {
		fmt.Println("Please Set Environ: REDIS for redis")
		fmt.Println("used like: export REDIS=\"123456@172.16.100.12:6379\"")
		return
	}
	addr, passwd := getRedisInfo()
	if len(addr) == 0 {
		fmt.Println("wrong format REDIS ", redisunc)
		return
	}

	mongounc := os.Getenv("MONGOUNC")
	if len(mongounc) == 0 {
		fmt.Println("Please Set Environ: MONGOUNC for mongo")
		fmt.Println("use like : export MONGOUNC=\"mongodb://172.16.100.12/golangtest?w=majority\"")
		return
	}
	if err := controllers.InitMongo(mongounc); err != nil {
		fmt.Println("mongo init error", err.Error())
	}

	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionProvider = "redis"
	beego.BConfig.WebConfig.Session.SessionProviderConfig = addr + ",10," + passwd
	beego.Run()
}
