package datastructure

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// SetClient set mongo client
func (m *MongoDB) SetClient(mongoClient *mongo.Client) {
	m.client = mongoClient
} // end of SetClient()

// TestConnect start ping and connect test
func (m *MongoDB) TestConnect(connTimeoutSec, pingTimeoutSec int) error {
	ctx, timeoutCancel := context.WithTimeout(context.Background(), time.Duration(connTimeoutSec)*time.Second)
	defer timeoutCancel()
	// connect test
	if err := m.client.Connect(ctx); err != nil {
		return err
	} //fi
	ctx, pingCancel := context.WithTimeout(context.Background(), time.Duration(pingTimeoutSec)*time.Second)
	defer pingCancel()
	if err := m.client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	return nil
} // TestConnect()

// SetDB set the mongo client focus DB by name
func (m *MongoDB) SetDB(DBName string) {
	m.database = DBName
}

// GaDBCollection get collction of DB
func (m *MongoDB) GaDBCollection(DBName, collectionName string) *mongo.Collection {
	return m.client.Database(DBName).Collection(collectionName)
}

// GaCollection GaDBCollection get collction of DB you set id with MongoDB.database
func (m *MongoDB) GaCollection(collectionName string) *mongo.Collection {
	return m.client.Database(m.database).Collection(collectionName)
}
