package dto

type WorkSpaceCreationDTO struct {
	Name        string `json:"name" validate:"required,min=3,max=80,printascii"`
	Slug        string `json:"slug" validate:"omitempty,min=3,max=100,lowercase,printascii"`
	Description string `json:"description" validate:"omitempty,max=1000"`
}