package database

import (
	"context"
	"log"
	"time"

	"ecom/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(cfg config.Config) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOpts := options.Client().
		ApplyURI(cfg.MongoURI).
		SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("mongo: connect error: %v", err)
	}

	return client
}

func PaymentCollection(client *mongo.Client, cfg config.Config) *mongo.Collection {
	col := client.Database(cfg.MongoDBName).Collection("payments")

	_, err := col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.M{"transaction_id": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Printf("mongo: create index error: %v", err)
	}

	return col
}

func ProductCollection(client *mongo.Client, cfg config.Config) *mongo.Collection {
	return client.Database(cfg.MongoDBName).Collection("products")
}

func TransactionCollection(client *mongo.Client, cfg config.Config) *mongo.Collection {
	return client.Database(cfg.MongoDBName).Collection("transactions")
}
