package dto

import (
	"DotConnect/schemas"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type WorkSpaceCreationDTO struct {
	Name        string `json:"name" validate:"required,min=3,max=80,printascii"`
	Slug        string `json:"slug" validate:"omitempty,min=3,max=100,lowercase,printascii"`
	Description string `json:"description" validate:"omitempty,max=1000"`
}


type AddMemberDTO struct{
	UserID string `json:"userId" validate:"required"`
    Role   string `json:"role"   validate:"required,oneof=admin member"`
}



type WorkSpaceResponseDTO struct {
	ID          bson.ObjectID     `json:"id"`
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Description string            `json:"description,omitempty"`
	Members     []schemas.WorkspaceMember `json:"members"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}