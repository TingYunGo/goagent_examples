package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"sync/atomic"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func init_handle(w http.ResponseWriter, r *http.Request) {
	db := get_db()
	var b []byte
	if db != nil {
		create_sql := "create table user2(id INTEGER PRIMARY KEY  AUTO_INCREMENT, name varchar(255), passwd varchar(255), age integer)"
		if dbvender == "sqlite3" {
			create_sql = "create table user2(id INTEGER PRIMARY KEY   AUTOINCREMENT, name varchar(255), passwd varchar(255), age integer);"
		}
		_, err := db.Exec(create_sql)
		if err != nil {
			b, _ = json.Marshal(map[string]interface{}{
				"status":  "error",
				"message": err.Error(),
			})
		} else {
			b, _ = json.Marshal(map[string]interface{}{
				"status":  "success",
				"message": "table user2 created.",
			})
		}
	}
	header := w.Header()
	header.Set("Cache-Control", "no-cache")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Connection", "Keep-Alive")
	header.Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func init() {
	fmt.Printf("main package init\n")
}

var global_db *sql.DB = nil
var dblock int32 = 0
var dbinited int32 = 0

func get_db() *sql.DB {
	if dbinited > 0 {
		return global_db
	}
	if atomic.AddInt32(&dblock, 1) == 1 {
		db, err := sql.Open(dbvender, dbunc)
		db.SetMaxOpenConns(400)
		if err == nil {
			global_db = db
			dbinited = 1
		} else {
			fmt.Println("Database Connect error:", err)
		}
	} else {
		atomic.AddInt32(&dblock, -1)
		for dbinited == 0 {
			time.Sleep(1)
		}
	}
	return global_db
}

func test1_handle(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	db := get_db()
	if db == nil {
		return
	}
	rows, err := db.Query("select id, name, age from user2")
	defer rows.Close()
	if err != nil {
		fmt.Println("select :", err)
		return
	}

	value := []interface{}{}
	for rows.Next() {
		var id int
		var name string
		var age int
		rows.Scan(&id, &name, &age)
		v := map[string]interface{}{
			"id":   id,
			"name": name,
			"age":  age,
		}
		value = append(value, v)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("other error:", err)
		return
	}
	b, _ := json.Marshal(map[string]interface{}{
		"status": "success",
		"result": "test1",
		"values": value,
	})
	header.Set("Cache-Control", "no-cache")
	header.Set("Connection", "Keep-Alive")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)

}

func db_transagent_update(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	db := get_db()
	if db == nil {
		return
	}

	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback()

	rs, err := tx.Exec("UPDATE user2 SET age=50 WHERE name='zhangsan'")
	if err != nil {
		fmt.Println(err)
	}
	rowAffected, err := rs.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rowAffected)

	rs, err = tx.Exec("UPDATE user2 SET age=15 WHERE name='lisi'")
	if err != nil {
		fmt.Println(err)
	}
	rowAffected, err = rs.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rowAffected)

	if err := tx.Commit(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("transagent success!")
	}

	b, _ := json.Marshal(map[string]interface{}{
		"status": "success",
		"result": "test1",
		"rows":   rowAffected,
	})
	header.Set("Cache-Control", "no-cache")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)

}
func query_one(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	db := get_db()
	if db == nil {
		return
	}
	id := 0
	name := ""
	age := 0
	err := db.QueryRow("select id, name, age from user2 order by id desc").Scan(&id, &name, &age)
	if err != nil {
		b, _ := json.Marshal(map[string]interface{}{
			"status":  "error",
			"result":  "query_one dbaccess error",
			"message": err.Error(),
		})
		header.Set("Cache-Control", "no-cache")
		header.Set("Connection", "Keep-Alive")
		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(b)

		return
	}

	b, _ := json.Marshal(map[string]interface{}{
		"status": "success",
		"result": "test1",
		"id":     id,
		"name":   name,
		"age":    age,
	})
	header.Set("Cache-Control", "no-cache")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)

}

func conn_query_one(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	db := get_db()
	if db == nil {
		return
	}
	id := 0
	name := ""
	age := 0
	conn, err := db.Conn(context.Background())
	if err != nil {
		fmt.Println("conn error: ", err)
		return
	}
	rows, err := conn.QueryContext(context.Background(), "select id, name, age from user2 order by id desc")
	if err != nil {
		fmt.Println("select :", err)
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
	fmt.Println("count :", cnt)

	b, _ := json.Marshal(map[string]interface{}{
		"status": "success",
		"result": "test1",
		"values": v,
	})
	header.Set("Cache-Control", "no-cache")
	header.Set("Connection", "Keep-Alive")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)

}

func InsertHandle(w http.ResponseWriter, r *http.Request) {
	db := get_db()
	if db == nil {
		return
	}
	insert_stmt := "INSERT INTO user2 (name, passwd, age) VALUES (?,?,?)"
	stmt, err := db.Prepare(insert_stmt)
	if stmt == nil {
		fmt.Println("sql prepare error ", err)
		return
	}
	defer stmt.Close()
	if err != nil {
		fmt.Println("db.Prepare error:", err)
		return
	}
	res, err := stmt.Exec("zhangsan", "abcdef", 17)
	if err != nil {
		fmt.Println("exec1 error: ", err)
		return
	}
	values := []interface{}{}
	id, err := res.LastInsertId()
	values = append(values, id)
	res, err = stmt.Exec("lisi", "abcdef", 18)
	if err != nil {
		fmt.Println("exec2 error: ", err)
		return
	}
	id, err = res.LastInsertId()
	values = append(values, id)

	b, _ := json.Marshal(map[string]interface{}{
		"status": "success",
		"result": "insert",
		"values": values,
	})
	header := w.Header()
	header.Set("Cache-Control", "no-cache")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Connection", "Keep-Alive")
	header.Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)

}

var dbunc = ""
var dbvender = "sqlite3"

func main() {
	listen := ":2000"
	dbunc = os.Getenv("DBUNC")
	if len(dbunc) == 0 {
		fmt.Println("Please Set Environ: DBUNC for database")
		return
	}
	if len(os.Args) > 1 {
		listen = os.Args[1]
	}
	fmt.Println("type UnsafePointer = ", (int)(reflect.UnsafePointer))
	http.HandleFunc("/init", init_handle)
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Cache-Control", "no-cache")
		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Connection", "Keep-Alive")
		header.Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(map[string]interface{}{
			"status": "success",
			"result": "test",
			"listen": listen,
		})
		w.Write(b)
	})
	http.HandleFunc("/test1", test1_handle)
	http.HandleFunc("/insert", InsertHandle)
	http.HandleFunc("/query_one", query_one)
	http.HandleFunc("/conn_query_one", conn_query_one)
	http.HandleFunc("/tx_update", db_transagent_update)
	fmt.Printf("listen %s\n", listen)

	http.ListenAndServe(listen, nil)
}
