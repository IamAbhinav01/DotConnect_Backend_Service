package repositories

import (
	"DotConnect/schemas"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ChannelRepository interface {
	CreateChannel(ctx context.Context,channelPayload *schemas.Channel) (*schemas.Channel, error)
	DuplicateChannelChecker(ctx context.Context,channelName string,workspaceId bson.ObjectID) (bool,error)
}

type ChannelRepositoryImpl struct{
	db *mongo.Collection
}

func NewChannelRepository(_db *mongo.Collection)*ChannelRepositoryImpl{
	return &ChannelRepositoryImpl{
		db: _db,
	}
}

func(repo *ChannelRepositoryImpl) CreateChannel(ctx context.Context,channelPayload *schemas.Channel) (*schemas.Channel, error){
	channelPayload.ID = bson.NewObjectID()
	channelPayload.CreatedAt = time.Now()
	channelPayload.UpdatedAt = time.Now()
	
	_,err := repo.db.InsertOne(ctx,channelPayload)
	if err!=nil{
		return nil,fmt.Errorf("%w",err)
	}

	return channelPayload,nil
}

func(repo *ChannelRepositoryImpl)DuplicateChannelChecker(ctx context.Context,channelName string,workspaceId bson.ObjectID) (bool,error){
	
	filter:=bson.M{
		"workspaceId":workspaceId,
		"name":channelName,
	}
	count,err := repo.db.CountDocuments(ctx,filter)
	if err != nil{
		return false,err
	}
	return count>0,err
}