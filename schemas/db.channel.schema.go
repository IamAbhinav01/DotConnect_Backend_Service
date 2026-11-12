package schemas

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)



type ChannelMember struct {
	UserID   string `bson:"userId" json:"userId"`
	AvatarURL string `bson:"avatar" json:"avatar"`
	Role     string `bson:"role" json:"role"`
	JoinedAt time.Time `bson:"joinedAt" json:"joinedAt"`
}

type Channel struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	WorkspaceID bson.ObjectID `bson:"workspaceId" json:"workspaceId"`
	Name string `bson:"name" json:"name"`
	Members []ChannelMember `bson:"members" json:"members"`
	CreatedAt   time.Time         `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time         `bson:"updatedAt" json:"updatedAt"`
}