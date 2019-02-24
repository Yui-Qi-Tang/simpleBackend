package datastructure

import (
	"github.com/mongodb/mongo-go-driver/mongo"
)

// MongoDB a struct for mongodb client manager
type MongoDB struct {
	client   *mongo.Client
	database string
}

// MongoDBSetting mongo db settings
type MongoDBSetting struct {
	Server string `yaml:"server"`
}
