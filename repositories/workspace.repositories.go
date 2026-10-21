package repositories

import "go.mongodb.org/mongo-driver/v2/mongo"

type WorkspaceRepository interface {
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

func (repo *WorkspaceRepositoryImpl) GetAllWorkspaces(){
	
}