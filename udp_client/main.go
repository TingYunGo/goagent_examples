package main

import (
	"fmt"
	"net"
	"time"
)

var id = 0

func main() {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 3000,
	})
	if err != nil {
		fmt.Println("udp bind error", err)
		return
	}
	defer socket.Close()

	go func() {
		for {
			data := [4096]byte{}
			if n, remoteAddr, err := socket.ReadFromUDP(data[:]); err == nil {
				fmt.Println("From:", remoteAddr, "Reply:", n, string(data[:n]))
			}
		}
	}()

	for {
		sendData := []byte(fmt.Sprintln("request from client", id))
		id++
		_, err = socket.Write(sendData) // 发送数据

		time.Sleep(time.Millisecond * 80)
	}
}
