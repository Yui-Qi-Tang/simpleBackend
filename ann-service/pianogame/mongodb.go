package pianogame

/*
   use mongo-go-driver
*/
import (
	"log"
	"simpleBackend/ann-service/pianogame/datastructure"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// Mongodb client for other package
var MongoGreeter datastructure.MongoDB

// init init. mongo db and return client
func init() {
	log.Println("init in pianogame mongodb")
	newMongoClient, err := mongo.NewClient("mongodb://localhost:27017") // mongodb from config file
	if err != nil {
		log.Fatalf("New client error: %v", err)
	} //fi

	MongoGreeter.SetClient(newMongoClient)
	if mErr := MongoGreeter.TestConnect(10, 2); mErr != nil {
		log.Fatalf("MongoGreeter error: %v", err)
	}
	log.Println("init in pianogame success")
} // end of initMongoDB
