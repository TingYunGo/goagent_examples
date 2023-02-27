## gopkg.in/mgo.v2框架嵌码示例代码

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
# DBUNC 格式参见 http://docs.mongodb.org/manual/reference/connection-string/
$ export DBUNC=172.16.100.12
$ ./mgo
```

### 访问测试
  使用浏览器或curl访问如下地址: <br/>
  添加测试数据: http://127.0.0.1:2002/insert <br/>
  查询单条数据: http://127.0.0.1:2002/query_one <br/>
  查询全部数据: http://127.0.0.1:2002/query_all <br/>
  查询Pipe数据: http://127.0.0.1:2002/pipe <br/>

  

### 应用性能数据查看
  在配置正确设置的情况下，应用程序启动后，登陆听云报表后台, 就能看到应用 mgo 的数据了。

