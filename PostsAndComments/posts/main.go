package main

import (
	"context"
	"database/sql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"posts/src"
	"time"
)

var (
	Client *mongo.Client
	DB     *sql.DB
)

func main() {
	logger := src.InitializeLogger(os.Stdout)
	mongoURI := "mongodb://localhost:27017" // замените на ваш URI
	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	Client = client

	psqlInfo := "host=localhost port=5432 user=postgres password=yourpassword dbname=yourdbname sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	DB = db

	if err := src.StartServer(logger); err != nil {
		logger.Fatalf("Server error: %v", err)
	}
}
