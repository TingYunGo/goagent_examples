## iris框架嵌码示例代码

源自 https://github.com/kataras/iris/tree/master/_examples/routing/basic
### 编译:
```bash
$ go mod tidy
$ go build
```
### 运行:
* 步骤1: 配置
  编辑tingyun.conf文件, 修改license_key 为实际的授权序列号, 修改collector.address 为实际的collector地址。

* 步骤2: 运行
  设置环境变量并运行:
```bash
$ export TINGYUN_GO_APP_CONFIG=tingyun.conf
$ ./iris_route_basic
```

### 访问测试
  使用浏览器或者使用curl访问如下任意链接: <br/>
  http://localhost:8080/  <br/>
  http://localhost:8080/home  <br/>
  http://localhost:8080/u/p:path  <br/>
  http://localhost:8080/u/username:string  <br/>
  http://localhost:8080/u/-1  <br/>
  http://localhost:8080/u/123  <br/>
  http://localhost:8080/u/firstname  <br/>
  http://localhost:8080/api/users/123  <br/>
  http://localhost:8080/admin  <br/>
  http://localhost:8080/admin/login  <br/>
  http://v1.localhost:8080  <br/>
  http://v1.localhost:8080/api/users/42  <br/>
  http://abc.localhost:8080/
  
### 应用性能数据查看
  在配置正确设置的情况下，应用程序启动后，登陆听云报表后台, 就能看到应用 iris_route_basic 的数据了。

