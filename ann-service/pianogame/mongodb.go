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
	var err error
	// Let Mongodb as global variable, you need to use '=' for assign instance for variable,
	// if use ':=' which is meat 'assign a new install for the variable' in this case, it's local variable in this function!!!
	Mongodb, err = mongo.NewClient("mongodb://localhost:27017") // mongodb from config file
	if err != nil {
		log.Fatalf("New client error: %v", err)
	} //fi

	conTimeOutCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err = Mongodb.Connect(conTimeOutCtx); err != nil {
		log.Fatalf("Client connection error: %v", err)
	} //fi

	pingTestCtx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	if err = Mongodb.Ping(pingTestCtx, readpref.Primary()); err != nil {
		log.Fatalf("Client ping mongodb server error: %v", err)
	} //fi
	log.Println("init in pianogame success")
} // end of initMongoDB

func gaCollection(DB string, collection string) *mongo.Collection {
	// bad idea
	return Mongodb.Database(DB).Collection(collection)
}
