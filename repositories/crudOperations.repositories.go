package repositories

import "go.mongodb.org/mongo-driver/v2/mongo"

type IcrudOperations interface {
	create()
	update()
	getById()
	getAll()
	delete()
}

type IcrudOperationsImpl struct {
	db *mongo.Client
}

func NewRepository( _db *mongo.Client) *IcrudOperationsImpl{
	return &IcrudOperationsImpl{
		db: _db,
	}
}

func (repo *IcrudOperationsImpl) create(){

}