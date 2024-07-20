package src

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

var (
	Client      *mongo.Client
	clentOnce   sync.Once
	clientErr   error
	isConnected bool
)

func InitMongoClient() error {
	clentOnce.Do(func() {
		ConnectingMongoDB()
	})
	return clientErr
}

func ConnectingMongoDB() {
	if isConnected {
		log.Print("Already connected to MongoDB")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Подключаемся к серверу MongoDB без указания конкретной базы данных
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		clientErr = err
		log.Print("Failed to connect to MongoDB:", err)
		return
	}

	err = Client.Ping(ctx, nil)
	if err != nil {
		clientErr = err
		log.Print("Failed to ping MongoDB:", err)
		return
	}

	isConnected = true
	log.Print("Connected to MongoDB server!")
}

func GetDatabase(dbName string) *mongo.Database {
	if err := EnsureMongoDBConnection(); err != nil {
		log.Printf("Failed to ensure MongoDB connection: %v", err)
		return nil
	}
	return Client.Database(dbName)
}

func GetCollection(dbName, collectionName string) *mongo.Collection {
	db := GetDatabase(dbName)
	if db == nil {
		return nil
	}
	return db.Collection(collectionName)
}

func EnsureMongoDBConnection() error {
	if isConnected && Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := Client.Ping(ctx, nil)
		if err == nil {
			return nil
		}
		isConnected = false
	}
	return InitMongoClient()
}
