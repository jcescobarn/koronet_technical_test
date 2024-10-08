package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConfig struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}

func NewMongoConfig(dbUser, dbPassword, dbName, dbHost, dbPort string) *MongoDBConfig {
	return &MongoDBConfig{
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		DBHost:     dbHost,
		DBPort:     dbPort,
	}
}

func (mc *MongoDBConfig) Connect() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin", mc.DBUser, mc.DBPassword, mc.DBHost, mc.DBPort, mc.DBName)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Println("MongoDB Connection Succesful")
	return client, nil
}
