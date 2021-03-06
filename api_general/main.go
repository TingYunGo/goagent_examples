package main

import (
	"errors"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	tingyun "github.com/TingYunGo/goagent"
)

var id int64 = 0

func PackHandler(data []byte, send func([]byte)) {
	fmt.Println("On pack in:", len(data))

	//听云API嵌码 -> 在协程中,事务开始时,
	//               创建一个事务追踪对象,用以追踪事务的执行性能
	action, _ := tingyun.CreateAction("UDP", "PackHandler")

	//听云API嵌码 -> 在协程中,事务结束时,
	//               调用tingyun.Action.Finish方法 结束事务追踪
	defer action.Finish()

	fmt.Println("ROUTINE:", tingyun.GetGID(), "Received: ", string(data))
	DoWorkerJob()

	currentID := atomic.AddInt64(&id, 1)
	routineID := tingyun.GetGID()
	ack := fmt.Sprintln("Routine:", routineID, "Anser From Server:", currentID)
	send([]byte(ack))
}

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

	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 3000,
	})
	if err != nil {
		fmt.Println("UDP 3000 Listen failed, err: ", err)
		return
	}
	defer listen.Close()
	fmt.Println("Started\n")
	for {
		var data [2048]byte
		if n, addr, err := listen.ReadFromUDP(data[:]); err == nil && n > 0 {
			go PackHandler(data[:n], func(ack []byte) {
				listen.WriteToUDP(ack, addr)
			})
		}
	}
}
