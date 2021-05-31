package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getConnection() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://webdev2:verynicewebdev2@cluster0.bzv65.mongodb.net/Webdev2?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()

	return client
}

var client *mongo.Client = getConnection()

func GetCollection(collectionName string) *mongo.Collection {
	return client.Database("Webdev2").Collection(collectionName)
}
