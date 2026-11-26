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

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ChannelService interface {
	CreateChannel(ctx context.Context, ownerID string, slug string, channelCreationDTO dto.ChannelCreationDTO) (*schemas.Channel, error)
	GetAllChannels(ctx context.Context, userId string, slug string) ([]schemas.Channel, error)
	FetchChannel(ctx context.Context, userID string, slug string, channelIDStr string) (*schemas.Channel, error)
	JoinChannel(ctx context.Context,userID string,slug string,channelIDStr string)(bool,error)
	AddMemberToChannel(ctx context.Context,userID string,slug string,channelIDStr string)(bool,error)
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

	is_auth := false
	for _, member := range workspace.Members {
		if member.UserID == ownerID && constants.Is_Manager(member.Role) {
			is_auth = true
		}
	}

	if !is_auth {
		return nil, fmt.Errorf("forbidden: only owners and admins can create channel")
	}

	exists, err := s.repo.DuplicateChannelChecker(ctx, name, workspace.ID)
	if err != nil {
		return nil, err
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
				UserID:    ownerID,
				Role:      constants.ROLE_ADMIN,
				AvatarURL: "/api/v1/users/" + ownerID + "/avatar",
				JoinedAt:  time.Now(),
			},
		},
	}

	createdChannel, err := s.repo.CreateChannel(ctx, channel)
	if err != nil {
		return nil, fmt.Errorf("create channel: %w", err)
	}
	return createdChannel, nil
}

func (s *ChannelServiceImpl) GetAllChannels(ctx context.Context, userId string, slug string) ([]schemas.Channel, error) {

	workspace, err := s.workspaceRepo.GetSlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("fetch workspace by slug %q: %w", slug, err)
	}

	is_allowed := false
	for _, member := range workspace.Members {
		if member.UserID == userId {
			is_allowed = true
		}
	}

	if !is_allowed {
		return nil, fmt.Errorf("forbidden: only members can view channels")
	}

	channels, err := s.repo.GetAllChannels(ctx, workspace.ID, userId)
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func (s *ChannelServiceImpl) FetchChannel(ctx context.Context, userID string, slug string, channelIDStr string) (*schemas.Channel, error) {

	workspace, err := s.workspaceRepo.GetSlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("fetch workspace by slug %q : %w", slug, err)
	}

	isMember := false
	isManager := false
	for _, member := range workspace.Members {
		if member.UserID == userID {
			isMember = true
			isManager = constants.Is_Manager(member.Role)
		}
	}

	if !isMember {
		return nil, fmt.Errorf("forbidden: only members can view channels")
	}

	channelID, err := bson.ObjectIDFromHex(channelIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid channel id %q", channelIDStr)
	}

	channel, err := s.repo.GetChannel(ctx, workspace.ID, channelID)
	if err != nil {
		return nil, err
	}

	// Private channels stay hidden for workspace members who are not in them.
	if channel.Public == nil || !*channel.Public {
		if !isManager {
			inChannel := false
			for _, member := range channel.Members {
				if member.UserID == userID {
					inChannel = true
				}
			}
			if !inChannel {
				return nil, fmt.Errorf("forbidden: you are not a member of this channel")
			}
		}
	}

	return channel, nil
}

func(s *ChannelServiceImpl) JoinChannel(ctx context.Context,userID string,slug string,channelIDStr string)(bool,error){

	workspace, err := s.workspaceRepo.GetSlug(ctx, slug)
	if err != nil {
		return false, fmt.Errorf("fetch workspace by slug %q : %w", slug, err)
	}

	is_member:=false

	for _,member:=range workspace.Members{
		if member.UserID == userID{
			is_member = true
			break
		}
	}

	if !is_member{
		return false, fmt.Errorf("forbidden: only workspace members can join channels")
	}


	channelID, err := bson.ObjectIDFromHex(channelIDStr)
	if err != nil {
		return false, fmt.Errorf("invalid channel id %q", channelIDStr)
	}

	newMember := schemas.ChannelMember{
		UserID: userID,
		Role: constants.ROLE_MEMBER,
		AvatarURL: "/api/v1/users/" + userID + "/avatar",
		JoinedAt: time.Now(),
	}

	ok,err:=s.repo.JoinChannel(ctx,workspace.ID,channelID,newMember)
	if err!=nil{
		return false,fmt.Errorf("Error occured while joining the channel")
	}
	if !ok{
		return false,fmt.Errorf("Operation failed while joining %w",err)
	}
	return true,nil
}

func(s *ChannelServiceImpl)AddMemberToChannel(ctx context.Context,requestUserID string,targetUserID string,slug string,channelIDStr string)(bool,error){

	workspace, err := s.workspaceRepo.GetSlug(ctx, slug)
	if err != nil {
		return false, fmt.Errorf("fetch workspace by slug %q : %w", slug, err)
	}

	isRequestedMember:=false
	isTargetedMember:=false
	isRequesterManager := false

	for _,member:=range workspace.Members{
		if member.UserID == requestUserID {
			isRequestedMember = true
			isRequesterManager = constants.Is_Manager(member.Role)
		}
		if member.UserID == targetUserID {
			isTargetedMember = true
		}
	}

	if !isRequestedMember {
		return false, fmt.Errorf("forbidden: requester is not a member of this workspace")
	}
	if !isTargetedMember {
		return false, fmt.Errorf("bad request: target user is not a member of this workspace")
	}

	channelID, err := bson.ObjectIDFromHex(channelIDStr)
	if err != nil {
		return false, fmt.Errorf("invalid channel id %q", channelIDStr)
	}

	channel, err := s.repo.GetChannel(ctx, workspace.ID, channelID)
	if err != nil {
		return false, fmt.Errorf("channel not found: %w", err)
	}

	isRequesterInChannel := false
	for _, m := range channel.Members {
		if m.UserID == requestUserID {
			isRequesterInChannel = true
			break
		}
	}

	if !isRequesterManager && !isRequesterInChannel {
		return false, fmt.Errorf("forbidden: you must be a member of this channel to add others")
	}



	newMember := schemas.ChannelMember{
		UserID: targetUserID,
		Role: constants.ROLE_MEMBER,
		AvatarURL: "/api/v1/users/" + targetUserID + "/avatar",
		JoinedAt: time.Now(),
	}

	ok, err := s.repo.AddMember(ctx, workspace.ID, channelID, newMember)
	if err != nil {
		return false, fmt.Errorf("error adding member to channel: %w", err) 
	}
	if !ok {
		return false, fmt.Errorf("operation failed while adding member")
	}
	return true, nil
}