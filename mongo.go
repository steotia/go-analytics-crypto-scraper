package main

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoDBClient() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
	// clientOptions := options.Client().ApplyURI("mongodb://localhost:27100")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		glog.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		glog.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	return client, err
}
