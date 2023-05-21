package storage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var Client *mongo.Client

func ConnectToMongo() {
	// Set client options
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%v:%v", os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT")))
	// Connect to MongoDB
	var err error
	Client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = Client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
}
