package main

import (
	// https://blog.golang.org/context
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo" // https://docs.mongodb.com/ecosystem/drivers/go/
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	// "github.com/mongodb/mongo-go-driver/bson"
)

// initMongoDB init. mongo db and return client
func initMongoDB() *mongo.Client {
	/*
		TO-DO:
			1. load config from file for mongodb
			    Base data: host/port/connect protocol
				Q: specify db and collection?
	*/
	client, err := mongo.NewClient("mongodb://localhost:27017") // 27017

	if err != nil {
		defer func() {
			log.Fatalf("read service config file error: %v", err)
		}()
	}
	return client
	/*
		if client, err := mongo.NewClient("mongodb://localhost:27017"); err != nil {
			fmt.Println("MongoDB connect failed!")
		} else {
			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
			fmt.Println(ctx)
			err = client.Connect(ctx) // ctx: connection timeout

			// mongo-go-driver: https://godoc.org/github.com/mongodb/mongo-go-driver/mongo

			// set DB and collection
			collection := client.Database("testing").Collection("numbers")

			ctx, _ = context.WithTimeout(context.Background(), 5*time.Second) // connection option
			// insert data to db.collection via bson
			res, err := collection.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159})

			if err != nil {
				fmt.Println("Insert Failed!!")
				return
			}
			id := res.InsertedID // get result if insert ok
			fmt.Println(id)
		} // fi
	*/
} // end of initMongoDB

func main() {
	/*
		TO-DO:
			1. load config from file for mongodb
			    Base data: host/port/connect protocol
				Q: specify db and collection?
			2. bind mongodb and go gin api together
	*/

	fmt.Println("Hello world, SimpleBackend!!")
	DBClient := initMongoDB()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	fmt.Println(ctx)
	err := DBClient.Connect(ctx) // ctx: connection timeout
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	if err = DBClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("ping error: %v", err)
	}
	DBClient.Disconnect(ctx)

	// Go Gin
	gin.SetMode(gin.TestMode) // enable server on localhost:8080
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// This handler will match /user/john but will not match /user/ or /user
	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	// Example login with JSON
	// Binding from JSON
	type Login struct {
		User     string `form:"user" json:"user" xml:"user"  binding:"required"`
		Password string `form:"password" json:"password" xml:"password" binding:"required"`
	}
	router.POST("/loginJSON", func(c *gin.Context) {
		var json Login
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if json.User != "manu" || json.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})
	// register account
	router.POST("user/register", func(c *gin.Context) {
		var registerData Login
		if err := c.ShouldBindJSON(&registerData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(registerData.User)
		c.JSON(http.StatusOK, gin.H{"status": "register ok!"})
	})
	router.Run() // listen and serve on 0.0.0.0:8080
}
