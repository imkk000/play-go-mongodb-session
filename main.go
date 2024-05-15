package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const (
	uri = "mongodb://localhost:27017/"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	opts := options.Client().
		ApplyURI(uri).
		SetReadPreference(readpref.PrimaryPreferred()).
		SetRetryWrites(true).
		SetWriteConcern(writeconcern.Majority())
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		fmt.Println("conn err:", err)
		return
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Println("ping err:", err)
		return
	}
	fmt.Println("Connected")

	col := client.Database("api").Collection("users")

	session, err := client.StartSession()
	if err != nil {
		fmt.Println("session err:", err)
		return
	}
	defer session.EndSession(ctx)
	fmt.Println("New Session")

	if _, err := session.WithTransaction(ctx, func(sc mongo.SessionContext) (_ any, err error) {
		_, err = col.InsertOne(sc, map[string]any{"name": "test"})
		return
	}); err != nil {
		fmt.Println("transaction err:", err)
		return
	}
	fmt.Println("Run TX")
}
