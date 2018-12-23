package main

import (
	"context" // https://blog.golang.org/context
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo" // https://docs.mongodb.com/ecosystem/drivers/go/
)

func main() {
	/*
		TO-DO:
			1. load config from file for mongodb
			    Base data: host/port/connect protocol
				Q: specify db and collection?
			2. bind mongodb and go gin api together
	*/
	fmt.Println("Hello world, SimpleBackend!!")
	if client, err := mongo.NewClient("mongodb://localhost:27017"); err != nil {
		fmt.Println("MongoDB connect failed!")
	} else {
		fmt.Println("OK!")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		fmt.Println(ctx)
		err = client.Connect(ctx) // ctx: connection timeout

		/*mongo-go-driver: https://godoc.org/github.com/mongodb/mongo-go-driver/mongo */

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
	}
}
