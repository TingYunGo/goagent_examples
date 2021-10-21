## beegov2框架嵌码示例代码

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
$ ./bee2
```

### 访问测试
  使用浏览器或者curl访问如下链接: <br/>
  http://localhost:8080/ <br/> 
  http://localhost:8080/api/123456 <br/> 
  http://localhost:8080/handler <br/> 
  http://localhost:8080/handler1 <br/> 
  http://localhost:8080/handler2
  

### 应用性能数据查看
  在配置正确设置的情况下，应用程序启动后，登陆听云报表后台, 就能看到应用 bee2 的数据了。

