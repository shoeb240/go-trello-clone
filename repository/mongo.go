package repository

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDBConnection() (*mongo.Database, error) {
	godotenv.Load(".env")
	connectionString := os.Getenv("CONNECTION_STRING")
	dbName := os.Getenv("DB_NAME")

	clientOption := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		panic("Could not connect to Mongodb")
	}

	DB := client.Database(dbName)

	return DB, nil
}
