package repositories

import (
	"DotConnect/schemas"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ChannelRepository interface {
	CreateChannel(ctx context.Context,channelPayload *schemas.Channel) (*schemas.Channel, error)
	DuplicateChannelChecker(ctx context.Context,channelName string,workspaceId bson.ObjectID) (bool,error)
	GetAllChannels(ctx context.Context,workspaceId bson.ObjectID,userID string)([]schemas.Channel,error)
	GetChannel(ctx context.Context,workspaceId bson.ObjectID,channelId bson.ObjectID)(*schemas.Channel,error)
	JoinChannel(ctx context.Context,workspaceId bson.ObjectID,channelID bson.ObjectID, member schemas.ChannelMember)(bool,error)
	AddMember(ctx context.Context,workspaceId bson.ObjectID,channelId bson.ObjectID,member schemas.ChannelMember)(bool,error)
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

func(repo *ChannelRepositoryImpl) GetAllChannels(ctx context.Context,workspaceId bson.ObjectID,userID string)([]schemas.Channel,error){
	
	filter:=bson.M{
		"workspaceId":workspaceId,
		"$or": []bson.M{
			{"public": true},
			{"members.userId": userID},
		},
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "createdAt", Value: 1}})

	cursor,err := repo.db.Find(ctx,filter,findOptions)
	if err != nil {
		return nil, fmt.Errorf("error querying channels: %w", err)
	}
	defer cursor.Close(ctx)

	channels := make([]schemas.Channel, 0) // if the channels return 0 json return empty array insread of nill
	if err := cursor.All(ctx, &channels); err != nil {
		return nil, fmt.Errorf("error decoding channels: %w", err)
	}
	return channels, nil
}

func(repo *ChannelRepositoryImpl)GetChannel(ctx context.Context,workspaceId bson.ObjectID,channelId bson.ObjectID)(*schemas.Channel,error){

	filter:=bson.M{
		"_id":         channelId,
		"workspaceId":workspaceId,
	}

	var channel schemas.Channel


	err:=repo.db.FindOne(ctx,filter).Decode(&channel)
	
	if err != nil {
		//document not found
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("channel not found")
		}
		return nil, fmt.Errorf("error fetching channel: %w", err)
	}
	return &channel,err
}

func(repo *ChannelRepositoryImpl)JoinChannel(ctx context.Context,workspaceId bson.ObjectID,channelID bson.ObjectID, member schemas.ChannelMember)(bool,error){


	filter:=bson.M{
		"_id":            channelID,
		"workspaceId":    workspaceId,
		"public":         true,
		"members.userId": bson.M{"$ne": member.UserID},
	}
	update:=bson.M{
		"$push":bson.M{"members":member},
		"$set":bson.M{"updatedAt":time.Now()},
	}

	result,err:=repo.db.UpdateOne(ctx,filter,update)
	if err != nil{
		return false,fmt.Errorf("error joining member: %w",err)
	}

	if result.MatchedCount == 0{
		return false,fmt.Errorf("channel not found , private, or you are already a member")
	}

	return true,nil

}

func(repo *ChannelRepositoryImpl)AddMember(ctx context.Context,workspaceId bson.ObjectID,channelId bson.ObjectID,member schemas.ChannelMember)(bool,error){

	filter:=bson.M{
		"_id":            channelId,
		"workspaceId":    workspaceId,
		"members.userId": bson.M{"$ne": member.UserID},
	}

	update:=bson.M{
		"$push": bson.M{
        "members": member, 
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}

	result,err:=repo.db.UpdateOne(ctx,filter,update)
	if err != nil{
		return false,fmt.Errorf("error adding the member: %w",err)
	}

	if result.MatchedCount == 0{
		return false,fmt.Errorf("channel not found , or you are already a member")
	}

	return true,nil
}