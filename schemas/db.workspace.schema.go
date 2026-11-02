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
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id"`
	InviteCode string        `bson:"inviteCode" json:"inviteCode"`
	UserID     string        `bson:"userId" json:"userId"`
	AvatarURL  string        `bson:"avatar" json:"avatar"`
	Status     string        `bson:"status" json:"status"` // "pending" | "approved" | "rejected"
	CreatedAt  time.Time     `bson:"createdAt" json:"createdAt"`
	ReviewedAt *time.Time    `bson:"reviewedAt,omitempty" json:"reviewedAt,omitempty"`
	ReviewedBy *string       `bson:"reviewedBy,omitempty" json:"reviewedBy,omitempty"`
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
	JoinRequests []JoinRequest     `bson:"joinRequests,omitempty" json:"joinRequests,omitempty"`
	CreatedAt   time.Time         `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time         `bson:"updatedAt" json:"updatedAt"`
}

