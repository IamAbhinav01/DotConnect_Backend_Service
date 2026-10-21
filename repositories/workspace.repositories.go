package repositories

import (
	"DotConnect/schemas"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type WorkspaceRepository interface {
	CreateWorkspace(ctx context.Context,workspacePayload *schemas.WorkSpace) (*schemas.WorkSpace,error)
	GetAllWorkspaces()
}

type WorkspaceRepositoryImpl struct {
	db *mongo.Collection
}

func NewWorkspaceRepository(_db *mongo.Collection) *WorkspaceRepositoryImpl{
	return  &WorkspaceRepositoryImpl{
		db: _db,
	}
}

func (repo *WorkspaceRepositoryImpl)CreateWorkspace(ctx context.Context,workspacePayload *schemas.WorkSpace) (*schemas.WorkSpace,error){

	workspacePayload.ID = bson.NewObjectID()
	workspacePayload.CreatedAt = time.Now()
	workspacePayload.UpdatedAt = time.Now()

	_,err:=repo.db.InsertOne(ctx,workspacePayload)

	if err != nil{
		return nil,fmt.Errorf("Error occured while creating the workspace")
	}
	return workspacePayload,nil
}



func (repo *WorkspaceRepositoryImpl) GetAllWorkspaces(){
	
}