package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/TingYunGo/libgo/ascii"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var databaseHandle *sql.DB = nil
var dblock int32 = 0
var dbinited int32 = 0

var dbunc = ""
var dbvender = "sqlite3"

//数据库
func getDB() *sql.DB {
	if dbinited > 0 {
		return databaseHandle
	}
	if atomic.AddInt32(&dblock, 1) == 1 {
		db, err := sql.Open(dbvender, dbunc)
		db.SetMaxOpenConns(400)
		if err == nil {
			databaseHandle = db
			dbinited = 1
		} else {
			fmt.Println("连接数据库失败:", err)
		}
	} else {
		atomic.AddInt32(&dblock, -1)
		for dbinited == 0 {
			time.Sleep(1)
		}
	}
	return databaseHandle
}

func getRedisInfo() (host, passwd string) {
	offset := ascii.StrRChr(redisunc, '@')
	if offset == -1 {
		return redisunc, ""
	}
	return redisunc[offset+1:], ascii.SubString(redisunc, 0, offset)
}
func countAccess(c *gin.Context) int {
	session := sessions.Default(c)
	count := 0
	if v := session.Get("count"); v != nil {
		count = v.(int)
	}
	session.Set("count", count+1)
	session.Save()
	return count
}
func answer(c *gin.Context, code int, b []byte) {
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "Keep-Alive")
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(code, string(b))
}

//HTTP route handlers:

func handlerQueryOne(c *gin.Context) {

	accessCount := countAccess(c)

	db := getDB()
	if db == nil {
		fmt.Println("database error.")
		c.String(http.StatusInternalServerError, "database access error")
		return
	}
	//使用DB的query方法遍历数据库数据
	id := 0
	name := ""
	age := 0
	err := db.QueryRow("select id, name, age from user2 order by id desc").Scan(&id, &name, &age)
	if err != nil {
		b, _ := json.Marshal(map[string]interface{}{
			"status":        "error",
			"result":        "query_one dbaccess error",
			"SessionAccess": accessCount,
			"message":       err.Error(),
		})
		answer(c, http.StatusOK, b)
		return
	}
	//如果有数据记录Next指针就不为true

	b, _ := json.Marshal(map[string]interface{}{
		"status":        "success",
		"result":        "test1",
		"SessionAccess": accessCount,
		"id":            id,
		"name":          name,
		"age":           age,
	})
	answer(c, http.StatusOK, b)
}
func handlerUpdate(c *gin.Context) {

	accessCount := countAccess(c)

	db := getDB()
	if db == nil {
		fmt.Println("database error.")
		c.String(http.StatusInternalServerError, "database access error")
		return
	}
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "database access error ", err.Error())
		return
	}
	defer tx.Rollback()

	rs, err := tx.Exec("UPDATE user2 SET age=50 WHERE name='zhangsan'")
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "UPDATE error ", err.Error())
		return
	}
	rowAffected, err := rs.RowsAffected()
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "RowsAffected error ", err.Error())
		return
	}
	rs, err = tx.Exec("UPDATE user2 SET age=15 WHERE name='lisi'")
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "UPDATE 2 error ", err.Error())
		return
	}
	rowAffected2, err := rs.RowsAffected()
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "RowsAffected error ", err.Error())
		return
	}
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "tx.Commit error ", err.Error())
		return
	}

	b, _ := json.Marshal(map[string]interface{}{
		"status":        "success",
		"result":        "update",
		"SessionAccess": accessCount,
		"rows":          rowAffected + rowAffected2,
	})
	answer(c, http.StatusOK, b)
}
func handlerQueryAll(c *gin.Context) {

	accessCount := countAccess(c)

	db := getDB()
	if db == nil {
		fmt.Println("database error.")
		c.String(http.StatusInternalServerError, "database access error")
		return
	}
	id := 0
	name := ""
	age := 0
	conn, err := db.Conn(context.Background())
	if err != nil {
		c.String(http.StatusInternalServerError, "database connect :"+err.Error())
		return
	}
	rows, err := conn.QueryContext(context.Background(), "select id, name, age from user2 order by id desc")
	if err != nil {
		fmt.Println("select:", err)
		c.String(http.StatusInternalServerError, "query error :"+err.Error())
		return
	}
	defer rows.Close()
	v := []interface{}{}
	cnt := 0
	for rows.Next() {
		cnt++
		if err := rows.Scan(&id, &name, &age); err == nil {
			v = append(v, map[string]interface{}{
				"id":   id,
				"name": name,
				"age":  age,
			})
		} else {
			fmt.Println("scan error: ", err)
		}
	}

	b, _ := json.Marshal(map[string]interface{}{
		"status":        "success",
		"result":        "test1",
		"SessionAccess": accessCount,
		"values":        v,
	})
	answer(c, http.StatusOK, b)
}
func handlerInsert(c *gin.Context) {

	accessCount := countAccess(c)

	db := getDB()
	if db == nil {
		fmt.Println("database error.")
		c.String(http.StatusInternalServerError, "database access error")
		return
	}
	insertStmt := "INSERT INTO user2 (name, passwd, age) VALUES (?,?,?)"
	stmt, err := db.Prepare(insertStmt)
	if stmt == nil {
		fmt.Println("sql prepare error ", err)
		c.String(http.StatusInternalServerError, "database access error : "+err.Error())
		return
	}
	defer stmt.Close()
	if err != nil {
		fmt.Println("db.Prepare 错误:", err)
		c.String(http.StatusInternalServerError, "database access error : "+err.Error())
		return
	}
	//如果有数据记录Next指针就不为true
	res, err := stmt.Exec("zhangsan", "abcdef", 17)
	if err != nil {
		fmt.Println("exec1 error: ", err)
		c.String(http.StatusInternalServerError, "sql execute error : "+err.Error())
		return
	}
	values := []interface{}{}
	id, err := res.LastInsertId()
	values = append(values, id)
	res, err = stmt.Exec("lisi", "abcdef", 18)
	if err != nil {
		fmt.Println("exec2 error: ", err)
		c.String(http.StatusInternalServerError, "sql execute error : "+err.Error())
		return
	}
	id, err = res.LastInsertId()
	values = append(values, id)
	resjson := map[string]interface{}{
		"status":        "success",
		"result":        "insert",
		"SessionAccess": accessCount,
		"values":        values,
	}
	if err != nil {
		resjson["warning"] = err.Error()
	}

	b, _ := json.Marshal(resjson)
	answer(c, http.StatusOK, b)
}
func handlerInit(c *gin.Context) {

	accessCount := countAccess(c)

	db := getDB()
	var b []byte
	if db == nil {
		fmt.Println("database error.")
		c.String(http.StatusInternalServerError, "database access error")
		return
	}

	createSQL := "create table user2(id INTEGER PRIMARY KEY   AUTOINCREMENT, name varchar(255), passwd varchar(255), age integer);"
	_, err := db.Exec(createSQL)
	if err != nil {
		b, _ = json.Marshal(map[string]interface{}{
			"status":        "error",
			"SessionAccess": accessCount,
			"message":       err.Error(),
		})
	} else {
		b, _ = json.Marshal(map[string]interface{}{
			"status":        "success",
			"SessionAccess": accessCount,
			"message":       "table user2 created.",
		})
	}
	answer(c, http.StatusOK, b)
}
func handlerWelcome(c *gin.Context) {

	firstname := c.DefaultQuery("firstname", "Guest")
	lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

	c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
}

var redisunc = ""

func main() {
	listen := ":2001"
	if len(os.Args) > 1 {
		listen = os.Args[1]
	}
	dbunc = os.Getenv("DBUNC")
	if len(dbunc) == 0 {
		fmt.Println("Please Set Environ: DBUNC for database")
		fmt.Println("used like: export DBUNC=\"./sqlite.db\"")
		return
	}

	redisunc = os.Getenv("REDIS")
	if len(redisunc) == 0 {
		fmt.Println("Please Set Environ: REDIS for redis")
		fmt.Println("used like: export REDIS=\"123456@172.16.100.12:6379\"")
		return
	}
	addr, passwd := getRedisInfo()
	if len(addr) == 0 {
		fmt.Println("wrong format REDIS ", redisunc)
		return
	}

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	router := gin.Default()

	store, _ := sessions.NewRedisStore(10, "tcp", addr, passwd, []byte("secret"))
	store.Options(sessions.Options{
		MaxAge:   60, //session 失效时间
		Secure:   true,
		HttpOnly: true,
	})
	router.Use(sessions.Sessions("my_session", store))

	router.GET("/query_one", handlerQueryOne)
	router.GET("/query_all", handlerQueryAll)
	router.GET("/update", handlerUpdate)
	router.GET("/init", handlerInit)
	router.GET("/insert", handlerInsert)
	router.GET("/welcome", handlerWelcome)
	fmt.Println("listen:", listen)
	router.Run(listen)
}
