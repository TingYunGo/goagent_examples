package main

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"encoding/json"
	"net/http"
	"os"
)

type Person struct {
	Name  string
	Phone string
}

func GetCollection(collectionName string) *mgo.Collection {
	if database == nil {
		return nil
	}
	return database.C(collectionName)
}

type Trainer struct {
	Name string
	Age  int
	City string
}

func query_onedata(conn *mgo.Collection) ([]byte, error) {
	filter := bson.M{"name": "Ale"}
	result := Person{}
	err := conn.Find(filter).One(&result)
	if err != nil {
		return nil, err
	}
	return json.Marshal(map[string]interface{}{
		"status": "success",
		"result": result,
	})
}
func query_onedata_none(conn *mgo.Collection) ([]byte, error) {

	result := Person{}
	err := conn.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return json.Marshal(map[string]interface{}{
		"status": "success",
		"result": result,
	})
}
func query_none(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	conn := GetCollection("people")
	if conn == nil {
		fmt.Println("No conn got")
	}
	b, err := query_onedata_none(conn)
	if err != nil {
		b, _ = json.Marshal(map[string]interface{}{
			"status": "error",
			"result": err.Error(),
		})
	}
	header.Set("Cache-Control", "no-cache")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Connection", "Keep-Alive")
	header.Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func query_one(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	conn := GetCollection("people")
	if conn == nil {
		fmt.Println("No conn got")
	}
	b, err := query_onedata(conn)
	if err != nil {
		b, _ = json.Marshal(map[string]interface{}{
			"status": "success",
			"result": err.Error(),
		})
	}
	header.Set("Cache-Control", "no-cache")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Connection", "Keep-Alive")
	header.Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
func query_alldata(conn *mgo.Collection) ([]byte, error) {

	var results []Person = nil
	if err := conn.Find(bson.D{{}}).All(&results); err != nil {
		return nil, err
	}

	return json.Marshal(map[string]interface{}{
		"status": "success",
		"count":  query_count(conn),
		"result": results,
	})

}

func query_pipe_data(conn *mgo.Collection) ([]byte, error) {

	var results []Person = nil

	m := []bson.M{
		{"$match": bson.M{}},
		{"$limit": 3},
	}

	if err := conn.Pipe(m).All(&results); err != nil {
		return nil, err
	}

	return json.Marshal(map[string]interface{}{
		"status": "success",
		"count":  query_count(conn),
		"result": results,
	})

}
func query_count(conn *mgo.Collection) int {
	cnt, err := conn.Count()
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return cnt
}
func query_all(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	conn := GetCollection("people")
	if conn == nil {
		fmt.Println("No conn got")
		return
	}
	b, err := query_alldata(conn)
	if err != nil {
		b, _ = json.Marshal(map[string]interface{}{
			"status": "success",
			"result": err.Error(),
		})
	}
	header.Set("Cache-Control", "no-cache")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Connection", "Keep-Alive")
	header.Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func query_pipe(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	conn := GetCollection("people")
	if conn == nil {
		fmt.Println("No conn got")
		return
	}
	b, err := query_pipe_data(conn)
	if err != nil {
		b, _ = json.Marshal(map[string]interface{}{
			"status": "success",
			"result": err.Error(),
		})
	}
	header.Set("Cache-Control", "no-cache")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Connection", "Keep-Alive")
	header.Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func InsertHandle(w http.ResponseWriter, r *http.Request) {
	conn := GetCollection("people")

	if conn == nil {
		return
	}

	err := conn.Insert(&Person{"Ale", "+55 53 8116 9639"},
		&Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		fmt.Println(err)
		return
	}

	b, _ := json.Marshal(map[string]interface{}{
		"status": "success",
		"result": "insert",
	})
	header := w.Header()
	header.Set("Cache-Control", "no-cache")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Connection", "Keep-Alive")
	header.Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

var session *mgo.Session = nil
var database *mgo.Database = nil

func main() {
	dbunc := os.Getenv("DBUNC")
	if len(dbunc) == 0 {
		fmt.Println("Please Set Environ: DBUNC for mongo")
		fmt.Println("DBUNC like : 172.16.100.12")
		return
	}
	ses, err := mgo.Dial(dbunc)
	if err != nil {
		fmt.Println(err)
		return
	}
	session = ses.Clone()
	ses.Close()
	database = session.DB("test")
	listen := ":2002"
	if len(os.Args) > 1 {
		listen = os.Args[1]
	}

	fmt.Println("Successfully connected and pinged.")

	http.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	http.HandleFunc("/query_one", query_one)
	http.HandleFunc("/query_none", query_none)
	http.HandleFunc("/query_all", query_all)
	http.HandleFunc("/pipe", query_pipe)
	http.HandleFunc("/insert", InsertHandle)
	fmt.Printf("listen %s\n", listen)

	http.ListenAndServe(listen, nil)
}
