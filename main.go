package main

import (
	// https://blog.golang.org/context
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo" // https://docs.mongodb.com/ecosystem/drivers/go/
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
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
		log.Fatalf("New client error: %v", err)
	} //fi

	conTimeOutCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err = client.Connect(conTimeOutCtx); err != nil {
		log.Fatalf("Client connection error: %v", err)
	} //fi

	pingTestCtx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	if err = client.Ping(pingTestCtx, readpref.Primary()); err != nil {
		log.Fatalf("Client ping mongodb server error: %v", err)
	} //fi
	log.Println("DB initial ok!")
	return client
} // end of initMongoDB

func gaCollection(c *mongo.Client, DB string, collection string) *mongo.Collection {
	return c.Database(DB).Collection(collection)
}

func main() {
	/*
		TO-DO:
			1. load config from file for mongodb
			    Base data: host/port/connect protocol
				Q: specify db and collection?
			2. bind mongodb and go gin api together
			3. refactor
	*/

	fmt.Println("Hello world, SimpleBackend!!")
	// set DB client
	DBClient := initMongoDB()

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

	router.POST("user/login", func(c *gin.Context) {
		var json Login
		// json decode
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// get db collection
		collection := gaCollection(DBClient, "testing", "user")
		// prepare filter to query
		filter := bson.M{
			"name":     json.User,
			"password": json.Password,
		}
		r := Login{}
		// query
		err := collection.FindOne(context.Background(), filter).Decode(&r)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Invalid username and password!"})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
		}
	})
	// register account
	router.POST("user/register", func(c *gin.Context) {
		var registerData Login
		if err := c.ShouldBindJSON(&registerData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		collection := gaCollection(DBClient, "testing", "user")
		filter := bson.M{
			"name":     registerData.User,
			"password": registerData.Password,
		}
		r := Login{}
		err := collection.FindOne(context.Background(), filter).Decode(&r)
		if err != nil {
			_, err := collection.InsertOne(context.Background(), filter)

			if err != nil {
				log.Fatalf("Insert one failed: %v", err)
			}
			c.JSON(http.StatusOK, gin.H{"status": "register ok!"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "This account has been registed!"})
		}
	})

	router.Run() // listen and serve on 127.0.0.1:8080 in gin.TestMode
}
