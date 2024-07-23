package src

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"
	"sync"
	"time"
)

var (
	client     *mongo.Client
	clientOnce sync.Once
	clientErr  error
	initDone   = make(chan struct{})
)

const mongoURI = "mongodb://localhost:27017"

func AsyncInsert(ctx context.Context, wg *sync.WaitGroup, collection *mongo.Collection, documents []interface{}) chan error {
	errChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()

		_, err := collection.InsertMany(ctx, documents)
		if err != nil {
			err = fmt.Errorf("failed to insert documents: %v", err)
		}
		errChan <- err
		close(errChan)
	}()

	return errChan
}

func AsyncRead(wg *sync.WaitGroup, collection *mongo.Collection, filter bson.D) chan []bson.M {
	resultChan := make(chan []bson.M)

	wg.Add(1)
	go func() {
		defer wg.Done()

		ctx := context.Background()

		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			resultChan <- nil
			close(resultChan)
			return
		}

		var results []bson.M
		if err = cursor.All(ctx, &results); err != nil {
			resultChan <- nil
			cursor.Close(ctx)
			close(resultChan)
			return
		}

		err = cursor.Close(ctx)
		if err != nil {
			log.Printf("Failed to close cursor: %v", err)
		}

		resultChan <- results
		close(resultChan)
	}()

	return resultChan
}

func Init(connectionString string) error {
	clientOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		rootCAs, err := x509.SystemCertPool()
		if err != nil {
			log.Print("Failed to create root CA pool:", err)
			clientErr = err
			close(initDone)
			return
		}
		tlsConfig := &tls.Config{
			RootCAs: rootCAs,
		}
		//релизовать проверкуу токена
		clientOptions := options.Client().ApplyURI(connectionString)
		clientOptions.SetTLSConfig(tlsConfig)
		client, clientErr = mongo.Connect(ctx, clientOptions)
		if clientErr != nil {
			log.Print("Failed to connect to MongoDB:", clientErr)
			close(initDone)
			return
		}
		clientErr = client.Ping(ctx, nil)
		if clientErr != nil {
			log.Print("Failed to ping MongoDB", clientErr)
			close(initDone)
			return
		}
		log.Print("Conncted to MongoDB server")

		err = client.Database("posts").RunCommand(ctx, nil).Err()
		if err != nil {
			log.Printf("Failed to create database 'posts': %v", err)
		}
	})
	<-initDone
	return clientErr
}
