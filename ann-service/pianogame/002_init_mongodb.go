package pianogame

/*
   use mongo-go-driver
*/
import (
	"log"
	"simpleBackend/ann-service/pianogame/datastructure"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// MongoGreeter client for other package
var MongoGreeter datastructure.MongoDB

// init init. mongo db and return client
func init() {
	newMongoClient, err := mongo.NewClient(MongoConfig.Server) // mongodb from config file
	if err != nil {
		log.Fatalf("New client error: %v", err)
	} //fi

	MongoGreeter.SetClient(newMongoClient)
	if mErr := MongoGreeter.TestConnect(10, 2); mErr != nil {
		log.Fatalf("MongoGreeter error: %v", err)
	}
	log.Println("init in mongo db success")
} // end of init
