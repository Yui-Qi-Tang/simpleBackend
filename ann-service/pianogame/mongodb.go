package pianogame

/*
   use mongo-go-driver
*/
import (
	"context"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
)

// Mongodb client for other package
var Mongodb *mongo.Client

// init init. mongo db and return client
func init() {
	log.Println("init in pianogame mongodb")
	/*
		TO-DO:
			1. load config from file for mongodb
			    Base data: host/port/connect protocol
				Q: specify db and collection?
	*/
	mongodb, err := mongo.NewClient("mongodb://localhost:27017") // mongodb from config file
	if err != nil {
		log.Fatalf("New client error: %v", err)
	} //fi

	conTimeOutCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err = mongodb.Connect(conTimeOutCtx); err != nil {
		log.Fatalf("Client connection error: %v", err)
	} //fi

	pingTestCtx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	if err = mongodb.Ping(pingTestCtx, readpref.Primary()); err != nil {
		log.Fatalf("Client ping mongodb server error: %v", err)
	} //fi
	Mongodb := mongodb
	log.Println(Mongodb)
} // end of initMongoDB

func gaCollection(DB string, collection string) *mongo.Collection {
	return Mongodb.Database(DB).Collection(collection)
}
