package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client

func GetCollection(uri, db, collectionName string) (*mongo.Collection, error) {
	var err error

	if dbClient != nil {
		return dbClient.Database(db).Collection(collectionName), nil
	}

	dbClient, err = connectDB(uri)

	if err != nil {
		return nil, err
	}

	return dbClient.Database(db).Collection(collectionName), nil
}

func connectDB(uri string) (*mongo.Client, error) {
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		return nil, err
	}

	fmt.Printf("Connected to instance: %s\n", uri)

	return client, nil
}
