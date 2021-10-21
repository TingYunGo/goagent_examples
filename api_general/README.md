## Go语言API嵌码示例代码
### 嵌码说明
  在这个UDP服务示例代码中, 服务处理框架和组件都是自定义的, 所以事务对象的创建和组件对象的创建全部采取手动嵌码方案。
### 代码嵌码对比:
![avatar](https://github.com/TingYunGo/goagent_examples/raw/main/api_general/api_general.jpg)
### 编译:
```bash
$ go mod tidy
$ go build
```
### 运行:
* 步骤1: 配置
  编辑tingyun.conf文件, 修改license_key 为实际的授权序列号, 修改collector.address 为实际的collector地址。

* 步骤2: 运行
  设置环境变量并运行
```bash
$ export TINGYUN_GO_APP_CONFIG=tingyun.conf
$ ./api_general &
```

### 访问测试
  使用本例相邻文件夹提供的udp_client发送数据测试。

```bash
$ cd ../udp_test && go build
$ ./udp_test
```

### 应用性能数据查看
  在嵌码无误和配置都正确设置的情况下，应用程序启动后，登陆听云报表后台,就能看到应用 api_general 的数据了。

