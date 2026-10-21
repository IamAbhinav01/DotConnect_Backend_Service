package dto

type WorkSpaceCreationDTO struct {
	Name        string `json:"name"  validate:"required,min=5,max=80,alphanum"`
	Slug        string `json:"slug" validate:"omitempty,min=3,max=100,lowercase,printascii"`
	Description string `json:"description" validate:"omitempty,max=1000"`
}