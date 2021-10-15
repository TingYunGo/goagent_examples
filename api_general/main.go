package main

import (
	"errors"
	"fmt"
	"net"
	"sync/atomic"
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

var id int64 = 0

func PackHandler(data []byte, send func([]byte)) {
	fmt.Println("On pack in:", len(data))
	action, _ := tingyun.CreateAction("UDP", "PackHandler")
	defer action.Finish()
	fmt.Println("ROUTINE:", tingyun.GetGID(), "Received: ", string(data))
	DoWorkerJob()
	currentID := atomic.AddInt64(&id, 1)
	send([]byte(fmt.Sprintln("Routine:", tingyun.GetGID(), "Anser From Server:", currentID)))
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
