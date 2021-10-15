## beego框架 + mongo + redis 嵌码示例代码

### 编译:
```bash
$ go mod tidy
$ go build
```
### 运行:
* 步骤1: 配置
  编辑tingyun.conf文件, 修改license_key 为实际的授权序列号, 修改collector.address 为实际的collector地址

* 步骤2: 运行
  设置环境变量并运行
```bash
$ export TINGYUN_GO_APP_CONFIG=tingyun.conf
```
  步骤3: 根据命令提示,设置redis 和 mongodb的地址
```bash
$ ./beego-mongodb-redis
```

### 访问测试
  使用浏览器或者curl访问如下链接: <br/>
  http://localhost:8080/insert <br/>
  http://localhost:8080/query_one 
  
  

### 应用性能数据查看
  登陆听云报表后台,查看应用 beego-mongodb-redis 数据

