package dto

type ChannelCreationDTO struct {
	Name        string `json:"name" validate:"required,min=3,max=80"`
	Public      *bool   `json:"public" validate:"omitempty"` 
}