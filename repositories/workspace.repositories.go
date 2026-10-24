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
	SlugAlreadyExists(ctx context.Context, slug string) (bool, error)
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
