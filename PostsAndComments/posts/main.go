package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"posts/src"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	src.Client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/posts.posts"))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = src.Client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	src.StartServer()
}
