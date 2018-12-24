package mongodb

/*
    use mongo-go-driver
*/
import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"log"
	"context"
	"time"
)


type mongoClient struct {
	Client *mongo.Client
}



// initMongoDB init. mongo db and return client
func InitMongoDB() mongoClient {
	/*
		TO-DO:
			1. load config from file for mongodb
			    Base data: host/port/connect protocol
				Q: specify db and collection?
	*/
	client, err := mongo.NewClient("mongodb://localhost:27017") // mongodb from config file
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
	
	var mongoC mongoClient
	mongoC.Client = client

	return mongoC
} // end of initMongoDB