package config

import (
	"context"
	"fmt"

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

	ClientOptions := options.Client().ApplyURI(MONGODB_URI)
	mongoClient,err := mongo.Connect(ClientOptions)

	if err!= nil{
		return fmt.Errorf("Error occured while connecting to mongodb")
	}


	err = mongoClient.Ping(context.Background(),nil)
	if err != nil{
		return fmt.Errorf("Error occured while pinging to the database")
	}

	Database := mongoClient.Database(DB_NAME)
	DotConnect = Database.Collection(DOTCONNECT_COLLECTION)

	return nil

}