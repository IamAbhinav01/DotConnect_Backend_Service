package repositories

import (
	"DotConnect/schemas"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type WorkspaceRepository interface {
	CreateWorkspace(ctx context.Context, workspacePayload *schemas.WorkSpace) (*schemas.WorkSpace, error)
	GetAllWorkspaces(ctx context.Context, userId string) ([]*schemas.WorkSpace, error)
	GetSlug(ctx context.Context,slug string)(*schemas.WorkSpace, error)
	SlugAlreadyExists(ctx context.Context, slug string) (bool, error)
	AddMemberToWorkspace(ctx context.Context,slug string,member schemas.WorkspaceMember)(bool,error)
	MemberExists(ctx context.Context,slug string,userId string) (bool, error)
	RemoveMemberFromWorkspace()
	UpdateWorkspace()
	PromoteMember(ctx context.Context,slug string,targetUser string, newRole string)(bool,error)
	DemoteMember(ctx context.Context,slug string,targetUser string, newRole string)(bool,error)
	DeleteWorkspace()
	GetAllMembersInWorkspace()
	JoinByInvite()
}

type WorkspaceRepositoryImpl struct {
	db *mongo.Collection
}

func NewWorkspaceRepository(_db *mongo.Collection) *WorkspaceRepositoryImpl {
	return &WorkspaceRepositoryImpl{
		db: _db,
	}
}

func (repo *WorkspaceRepositoryImpl) CreateWorkspace(ctx context.Context, workspacePayload *schemas.WorkSpace) (*schemas.WorkSpace, error) {

	workspacePayload.ID = bson.NewObjectID()
	workspacePayload.CreatedAt = time.Now()
	workspacePayload.UpdatedAt = time.Now()

	_, err := repo.db.InsertOne(ctx, workspacePayload)

	if err != nil {
		return nil, fmt.Errorf("error occurred while creating the workspace: %w", err)
	}
	return workspacePayload, nil
}

func (repo *WorkspaceRepositoryImpl) GetAllWorkspaces(ctx context.Context, userId string) ([]*schemas.WorkSpace, error) {

	cursor, err := repo.db.Find(ctx, bson.M{"members.userId": userId}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		return nil, fmt.Errorf("error fetching workspaces: %w", err)
	}
	defer cursor.Close(ctx)

	var workspaces []*schemas.WorkSpace

	err = cursor.All(ctx, &workspaces)
	if err != nil {
		return nil, fmt.Errorf("error decoding workspaces: %w", err)
	}
	return workspaces, nil
}

func (repo *WorkspaceRepositoryImpl) SlugAlreadyExists(ctx context.Context, slug string) (bool, error) {
	count, err := repo.db.CountDocuments(ctx, bson.M{"slug": slug})
	return count > 0, err
}

func (repo *WorkspaceRepositoryImpl)MemberExists(ctx context.Context,slug string,userId string) (bool, error){
	
	count, err := repo.db.CountDocuments(ctx, bson.M{"slug": slug,"members.userId":userId})
	if err != nil {
		return false, err
	}

	return count > 0, err
}

func(repo *WorkspaceRepositoryImpl) AddMemberToWorkspace(ctx context.Context,slug string,member schemas.WorkspaceMember)(bool,error){
	
	filter:= bson.M{"slug":slug}
	update:=bson.M{
		"$addToSet":bson.M{"members":member},
		"$set":bson.M{"JoinedAt":time.Now()},
	}


	result,err := repo.db.UpdateOne(ctx,filter,update)
	if err != nil{
		return false,fmt.Errorf("error adding member: %w",err)
	}

	if result.MatchedCount == 0{
		return false,fmt.Errorf("workspace not found")
	}

	return true,nil
}

func (repo *WorkspaceRepositoryImpl) GetSlug(ctx context.Context,slug string)(*schemas.WorkSpace, error){

	var workspaces *schemas.WorkSpace
	err := repo.db.FindOne(ctx,bson.M{"slug":slug}).Decode(&workspaces)
	if err != nil{
		return nil, err
	}

	return workspaces,nil
}
func (repo *WorkspaceRepositoryImpl)PromoteMember(ctx context.Context,slug string,targetUser string, newRole string)(bool,error){

	filter:=bson.M{"slug":slug,"members.userId":targetUser}
	update:=bson.M{
		"$set":bson.M{"members.$.role":newRole,
			"updatedAt":time.Now()},
	}

	result,err := repo.db.UpdateOne(ctx,filter,update)
	if err != nil{
		return false,fmt.Errorf("error promoting member: %w",err)
	}

	if result.MatchedCount == 0{
		return false,fmt.Errorf("workspace not found or user is not a member of this workspace")
	}

	return true,nil
}
func (repo *WorkspaceRepositoryImpl)DemoteMember(ctx context.Context,slug string,targetUser string, newRole string)(bool,error){

	filter:=bson.M{"slug":slug,"members.userId":targetUser}
	update:=bson.M{
		"$set":bson.M{"members.$.role":newRole,
			"updatedAt":time.Now()},
	}

	result,err := repo.db.UpdateOne(ctx,filter,update)
	if err != nil{
		return false,fmt.Errorf("error demoting member: %w",err)
	}

	if result.MatchedCount == 0{
		return false,fmt.Errorf("workspace not found or user is not a member of this workspace")
	}

	return true,nil
}