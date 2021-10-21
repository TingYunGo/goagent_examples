## 数据库框架嵌码示例代码

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
$ export DBUNC=./sqlite.db
$ ./dbexample_sqlite
```

### 访问测试
  使用浏览器或curl访问如下地址: <br/>
  数据库初始化: http://127.0.0.1:2000/init <br/>
  添加测试数据: http://127.0.0.1:2000/insert <br/>
  查询测试数据: http://127.0.0.1:2000/query_one
  

### 应用性能数据查看
  在配置正确设置的情况下，应用程序启动后，登陆听云报表后台, 就能看到应用 dbexample_sqlite 的数据了。

