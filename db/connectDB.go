package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var TaskCollection *mongo.Collection
var Context = context.Background()

func ConnectDB(connectionString string) {
	MongoUri := connectionString
	clientOptions := options.Client().ApplyURI(MongoUri)
	client, err := mongo.Connect(Context, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(Context, nil)
	if err != nil {
		log.Fatal(err)
	}
	TaskCollection = client.Database("tasker").Collection("tasks")
	fmt.Println("DB connection success")
}
