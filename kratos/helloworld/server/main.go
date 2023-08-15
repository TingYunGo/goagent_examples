package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"

	"helloworld/helloworld"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "helloworld"
	// Version is the version of the compiled software.
	// Version = "v1.0.0"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	helloworld.UnimplementedGreeterServer
}

/*
func Callers(skip int, pc []uintptr) int {
	// runtime.callers uses pc.array==nil as a signal
	// to print a stack trace. Pick off 0-length pc here
	// so that we don't let a nil pc slice get to it.
	if len(pc) == 0 {
		return 0
	}
	return callers(skip, pc)
}

*/
func getnameByAddr(p interface{}) string {
	ptr, _ := strconv.ParseInt(fmt.Sprintf("%x", p), 16, 64)
	return runtime.FuncForPC(uintptr(ptr)).Name()
}
func callStack(skip int) []string {
	var slice []string
	slice = make([]string, 0, 15)
	opc := uintptr(0)
	for i := skip + 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if opc == pc {
			continue
		}
		fname := getnameByAddr(pc)
		index := strings.Index(fname, "/tingyun/")
		if index > 0 {
			continue
		}
		opc = pc
		//截断源文件名
		index = strings.Index(file, "/src/")
		if index > 0 {
			file = file[index+5 : len(file)]
		}
		slice = append(slice, fmt.Sprintf("%s(%s:%d)", fname, file, line))
	}
	return slice
}
func PrintCaller() {
	fmt.Println(callStack(1))
	return
	stackList := make([]uintptr, 30)
	count := runtime.Callers(2, stackList)
	for i := 0; i < count; i++ {
		name := runtime.FuncForPC(stackList[i]).Name()
		fmt.Println(i, name)
	}
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	fmt.Println("request in")
	if in.Name == "error" {
		fmt.Println("name error\n")
		return nil, errors.BadRequest("custom_error", fmt.Sprintf("invalid argument %s", in.Name))
	}
	PrintCaller()
	if in.Name == "panic" {
		panic("server panic")
	}
	return &helloworld.HelloReply{Message: fmt.Sprintf("Hello %+v", in.Name)}, nil
}

func main() {
	s := &server{}
	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			recovery.Recovery(),
		),
	)
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			recovery.Recovery(),
		),
	)
	helloworld.RegisterGreeterServer(grpcSrv, s)
	helloworld.RegisterGreeterHTTPServer(httpSrv, s)

	app := kratos.New(
		kratos.Name(Name),
		kratos.Server(
			httpSrv,
			grpcSrv,
		),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
