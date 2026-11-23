package services

import (
	"DotConnect/dto"
	"DotConnect/repositories"
	"DotConnect/schemas"
	"DotConnect/utils/constants"
	"context"
	"fmt"
	"strings"
	"time"
)

type ChannelService interface {
	CreateChannel(ctx context.Context, ownerID string, slug string, channelCreationDTO dto.ChannelCreationDTO) (*schemas.Channel, error)
	GetAllChannels(ctx context.Context,userId string,slug string)([]schemas.Channel,error)
}

type ChannelServiceImpl struct {
	repo          repositories.ChannelRepository
	workspaceRepo repositories.WorkspaceRepository
}

func NewChannelService(_repo repositories.ChannelRepository, _workspaceRepo repositories.WorkspaceRepository) *ChannelServiceImpl {
	return &ChannelServiceImpl{
		repo:          _repo,
		workspaceRepo: _workspaceRepo,
	}
}

func (s *ChannelServiceImpl) CreateChannel(ctx context.Context, ownerID string, slug string, channelCreationDTO dto.ChannelCreationDTO) (*schemas.Channel, error) {
	name := strings.TrimSpace(channelCreationDTO.Name)
	if name == "" {
		return nil, fmt.Errorf("channel name cannot be empty")
	}

	workspace, err := s.workspaceRepo.GetSlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("fetch workspace by slug %q: %w", slug, err)
	}

	is_auth:=false
	for _,member:=range workspace.Members{
		if member.UserID == ownerID && constants.Is_Manager(member.Role){
			is_auth = true
		} 
	}

	if !is_auth{
		return nil,fmt.Errorf("forbidden: only owners and admins can create channel")
	}


	exists,err:=s.repo.DuplicateChannelChecker(ctx,name,workspace.ID)
	if err!=nil{
		return nil,err
	}

	if exists {
    	return nil, fmt.Errorf("channel with name %q already exists in this workspace", name)
	}

	channel := &schemas.Channel{
		WorkspaceID: workspace.ID,
		Name:        name,
		Public:      channelCreationDTO.Public,
		CreatedBy:   ownerID,
		Members: []schemas.ChannelMember{
			{
				UserID: ownerID,
				Role: constants.ROLE_ADMIN,
				AvatarURL:  "/api/v1/users/" + ownerID + "/avatar",
				JoinedAt: time.Now(),
			},
		},
	}

	createdChannel, err := s.repo.CreateChannel(ctx, channel)
	if err != nil {
		return nil, fmt.Errorf("create channel: %w", err)
	}
	return createdChannel, nil
}

func (s *ChannelServiceImpl) GetAllChannels(ctx context.Context,userId string,slug string)([]schemas.Channel,error){
	
	workspace, err := s.workspaceRepo.GetSlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("fetch workspace by slug %q: %w", slug, err)
	}

	is_allowed := false
	for _,member := range workspace.Members{
		if member.UserID == userId{
			is_allowed = true
		}
	}

	if !is_allowed{
		return nil,fmt.Errorf("forbidden: only members can view channels")
	}

	channels,err := s.repo.GetAllChannels(ctx,workspace.ID,userId)
	if err!= nil{
		return nil,err
	}
	return channels,nil
}