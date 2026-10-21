package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	Client *mongo.Client
	Db *mongo.Database
	DotConnect *mongo.Collection
)

func ConnectToDB() error{
	MONGODB_URI := GetString("MONGODB_URI")
	DB_NAME := GetString("DB_NAME")
	DOTCONNECT_COLLECTION := GetString("DOTCONNECT_COLLECTION")

	if MONGODB_URI == ""{
		return fmt.Errorf("Error occured while decoding the MONGODB URI")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	ClientOptions := options.Client().ApplyURI(MONGODB_URI).SetTLSConfig(tlsConfig)
	mongoClient, err := mongo.Connect(ClientOptions)
	if err != nil {
		return fmt.Errorf("Error occured while connecting to mongodb: %w", err)
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("Error occured while pinging to the database: %w", err)
	}

	Client = mongoClient
	Db = mongoClient.Database(DB_NAME)
	DotConnect = Db.Collection(DOTCONNECT_COLLECTION)
	if err := ensureIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}
	return nil

}


func ensureIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		// 1. Enforce globally unique slug (e.g. "abhinav-ws")
		{
			Keys:    bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		// 2. Fast querying for "All workspaces where User X is a member"
		{
			Keys: bson.D{{Key: "members.userId", Value: 1}},
		},
	}
	_, err := DotConnect.Indexes().CreateMany(ctx, indexModels)
	return err
}