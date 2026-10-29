package schemas

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)


type WorkspaceInvite struct{
	ID bson.ObjectID `bson:"id" json:"id"`
	Code string `bson:"code" json:"code"`
	CreatedBy string `bson:"createdby" json:"createdby"`
	MaxApprovedJoins int           `bson:"maxApprovedJoins" json:"maxApprovedJoins"`
	ApprovedCount    int           `bson:"approvedCount" json:"approvedCount"`  
	CreatedAt   time.Time         `bson:"createdAt" json:"createdAt"`
	ExpiresAt   time.Time         `bson:"expiresAt" json:"expiresAt"`
}


type JoinRequest struct{

}

type WorkspaceMember struct {
	UserID   string `bson:"userId" json:"userId"`
	AvatarURL string `bson:"avatar" json:"avatar"`
	Role     string `bson:"role" json:"role"`
	JoinedAt time.Time `bson:"joinedAt" json:"joinedAt"`
}

type WorkSpace struct{
	ID bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name" json:"name"`
	Slug        string            `bson:"slug" json:"slug"`   
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	Members []WorkspaceMember `bson:"members" json:"members"`
	Invites []WorkspaceInvite `bson:"invites,omitempty" json:"invites,omitempty"`
	CreatedAt   time.Time         `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time         `bson:"updatedAt" json:"updatedAt"`
}

