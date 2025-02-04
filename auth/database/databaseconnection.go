package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client{
	err := godotenv.Load(".env")
	if err!=nil{
		log.Fatal("Error loading .env file")
	}

	MongoDB := os.Getenv("MONGODB_URI")

	client,err := mongo.NewClient(options.Client().ApplyURI(MongoDB))
	if err!=nil{
		log.Fatal(err)
	}

	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	connectionErr := client.Connect(ctx)
	if connectionErr!=nil{
		log.Fatal(connectionErr)
	}

	fmt.Println("Connected to mongodb")

	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection{
	var collection *mongo.Collection = client.Database("cluster0").Collection(collectionName)
	return collection
}