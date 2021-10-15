package controllers

import (
	"context"
	"fmt"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MainController struct {
	beego.Controller
}

var mongounc = ""
var client *mongo.Client = nil

func GetCollection() *mongo.Collection {
	if client == nil {
		return nil
	}
	return client.Database("golangtest").Collection("my_collection")
}

func InitMongo(unc string) error {

	mongounc = unc
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(mongounc).SetMaxPoolSize(20))
	if err != nil {
		return err
	}
	client = c
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		c.Disconnect(ctx)
		cancel()
		return err
	}
	fmt.Println("Successfully connected and pinged.")
	return nil
}

func (c *MainController) sessionCountAccess() int {
	count := 0
	if v := c.GetSession("Count"); v != nil {
		count = v.(int)
	}
	c.SetSession("Count", count+1)
	return count
}
func (c *MainController) Get() {
	switch c.Ctx.Request.RequestURI {
	case "/insert":
		c.handleFunction(insert_data)
		return
	case "/query_one":
		c.handleFunction(query_onedata)
		return
	case "/delete_many":
		c.handleFunction(delete_many)
		return
	case "/update_many":
		c.handleFunction(update_many)
		return
	case "/query_none":
		c.handleFunction(query_none)
		return
	case "/query_all":
		c.handleFunction(query_alldata)
		return
	default:
		c.handleDefault()
		return
	}
}
func (c *MainController) handleDefault() {
	accessCount := c.sessionCountAccess()
	c.Data["Website"] = "beego.me"
	c.Data["RequestURI"] = c.Ctx.Request.RequestURI
	c.Data["Email"] = "astaxie@gmail.com"
	c.Data["SessionAccess"] = accessCount
	c.TplName = "index.tpl"
}
func (c *MainController) handleFunction(handler func(*mongo.Collection) map[string]interface{}) {
	conn := GetCollection()
	if conn == nil {
		c.Data["json"] = map[string]interface{}{
			"status":        "error",
			"SessionAccess": c.sessionCountAccess(),
			"error":         "Collection is nil",
		}
	} else {
		r := handler(conn)
		r["SessionAccess"] = c.sessionCountAccess()
		c.Data["json"] = r
	}
	c.ServeJSON()
}

type Trainer struct {
	Name       string
	Age        int
	City       string
	InsertTime string
}

func query_onedata(conn *mongo.Collection) map[string]interface{} {
	filter := bson.D{{"name", "Ash"}}
	var result Trainer
	err := conn.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		}
	}
	return map[string]interface{}{
		"status": "success",
		"result": result,
	}
}
func delete_many(conn *mongo.Collection) map[string]interface{} {
	filter := bson.D{{"name", "Ash"}}
	var result Trainer
	res, err := conn.DeleteMany(context.TODO(), filter)
	if err != nil {
		return map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		}
	}
	return map[string]interface{}{
		"status":      "success",
		"deleteCount": res.DeletedCount,
		"result":      result,
	}
}
func update_many(conn *mongo.Collection) map[string]interface{} {
	filter := bson.D{{"name", "Ash"}}
	var result Trainer
	ash := Trainer{"Bsh", 16, "Beijing", time.Now().Format("2006-01-02 15:04:05.000")}
	update := bson.D{{"$set",
		bson.D{
			{"name", ash.Name},
			{"age", ash.Age},
			{"city", ash.City},
			{"insertTime", ash.InsertTime},
		},
	}}
	res, err := conn.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		}
	}
	return map[string]interface{}{
		"status":        "success",
		"modifiedCount": res.ModifiedCount,
		"result":        result,
	}
}
func insert_data(conn *mongo.Collection) map[string]interface{} {
	now := time.Now().Format("2006-01-02 15:04:05.000")
	ash := Trainer{"Ash", 10, "Pallet Town", now}
	// misty := Trainer{"Misty", 10, "Cerulean City"}
	// brock := Trainer{"Brock", 15, "Pewter City"}
	//插入某一条数据
	insertID := ""
	iResult, err := conn.InsertOne(context.TODO(), ash)
	if err == nil {
		id := iResult.InsertedID.(primitive.ObjectID)
		insertID = id.Hex()
	}
	if err != nil {
		return map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		}
	}
	return map[string]interface{}{
		"status":   "success",
		"insertID": insertID,
		"result":   "",
	}
}
func query_alldata(conn *mongo.Collection) map[string]interface{} {

	findOptions := options.Find()
	findOptions.SetLimit(20)

	var results []*Trainer

	cur, err := conn.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		}
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var elem Trainer
		err := cur.Decode(&elem)
		if err != nil {
			return map[string]interface{}{
				"status": "error",
				"error":  err.Error(),
			}
		}
		results = append(results, &elem)
	}
	return map[string]interface{}{
		"status": "success",
		"Count":  len(results),
		"result": results,
	}
}

func query_none(conn *mongo.Collection) map[string]interface{} {
	filter := bson.D{{"name1", "111"}}
	var result Trainer
	err := conn.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return map[string]interface{}{
			"status": "error",
			"result": err.Error(),
		}
	}
	return map[string]interface{}{
		"status": "success",
		"result": result,
	}
}
