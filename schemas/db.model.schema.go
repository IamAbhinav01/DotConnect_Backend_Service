package schemas

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const RoleOwner string = "owner"
const RoleAdmin string = "admin"
const RoleMemebr string = "member"

type WorkspaceMember struct {
	UserID   string `bson:"userId" json:"userId"`
	Role     string `bson:"role" json:"role"`
	JoinedAt time.Time `bson:"joinedAt" json:"joinedAt"`
}

type WorkSpace struct{
	ID bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name" json:"name"`
	Slug        string            `bson:"slug" json:"slug"`   
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	Members []WorkspaceMember `bson:"members" json:"members"`
	CreatedAt   time.Time         `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time         `bson:"updatedAt" json:"updatedAt"`
}

type WorkSpaceCreationDTO struct{
	Name string `bson:"name" json:"name"`
	Slug        string `json:"slug,omitempty"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
}