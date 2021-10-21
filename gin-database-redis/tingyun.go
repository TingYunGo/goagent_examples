package main

import (
	// 示例代码中使用了数据库
	// 这行代码的加入实现了对数据库操作的性能数据采集.
	_ "github.com/TingYunGo/goagent/database"

	// 示例代码中使用了gin框架
	// 这行代码的加入,实现了对gin框架的自动嵌码
	_ "github.com/TingYunGo/goagent/frameworks/gin"

	// 示例代码中使用了redigo库访问redis
	// 这行代码的加入,实现了对redis操作的性能数据采集.
	_ "github.com/TingYunGo/goagent/nosql/redigo"
)
