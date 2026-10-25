package services

import (
	"DotConnect/dto"
	"DotConnect/repositories"
	"DotConnect/schemas"
	"DotConnect/utils/generator"
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type WorkspaceService interface {
	CreateWorkspace(ctx context.Context, ownerId string, dto dto.WorkSpaceCreationDTO) (*schemas.WorkSpace, error)
	GetUserWorkspace(ctx context.Context, userId string) ([]dto.WorkSpaceResponseDTO, error)
	AddMemberToWorkspace(ctx context.Context,slug string,requsterId string,member dto.AddMemberDTO) (bool,error)
}

type WorkspaceServiceImpl struct {
	repo repositories.WorkspaceRepository
}

func NewWorkspaceService(_repo repositories.WorkspaceRepository) *WorkspaceServiceImpl {
	return &WorkspaceServiceImpl{
		repo: _repo,
	}
}

func (s *WorkspaceServiceImpl) CreateWorkspace(ctx context.Context, ownerId string, dto dto.WorkSpaceCreationDTO) (*schemas.WorkSpace, error) {

	name := strings.TrimSpace(dto.Name)
	if name == "" {
		return nil, fmt.Errorf("workspace name cannot be empty")
	}
	description := dto.Description
	slug := dto.Slug

	if strings.TrimSpace(slug) == "" {
		slug = generator.Slugify(dto.Name)
	}

	baseSlug := slug
	for {
		exists, err := s.repo.SlugAlreadyExists(ctx, slug)
		if err != nil {
			return nil, fmt.Errorf("check workspace slug: %w", err)
		}
		if exists {
			slug = generator.SlugCollisionHandler(baseSlug)
			continue
		}

		ws := &schemas.WorkSpace{
			Name:        name,
			Description: description,
			Slug:        slug,
			Members: []schemas.WorkspaceMember{
				{
					UserID:   ownerId,
					AvatarURL: "/api/v1/users/" + ownerId + "/avatar",
					Role:     schemas.RoleOwner,
					JoinedAt: time.Now(),
				},
			},
		}

		createdWorkspace, err := s.repo.CreateWorkspace(ctx, ws)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				slug = generator.SlugCollisionHandler(baseSlug)
				continue
			}
			return nil, err
		}
		return createdWorkspace, nil
	}
}
func (s *WorkspaceServiceImpl) GetUserWorkspace(ctx context.Context, userId string) ([]dto.WorkSpaceResponseDTO, error) {
	workspaces, err := s.repo.GetAllWorkspaces(ctx, userId)
	if err != nil {
		return nil, err
	}

	result := make([]dto.WorkSpaceResponseDTO, len(workspaces))
	for i, workspace := range workspaces {
		result[i] = dto.WorkSpaceResponseDTO{
			ID:          workspace.ID,
			Name:        workspace.Name,
			Slug:        workspace.Slug,
			Description: workspace.Description,
			Members:     workspace.Members,
			CreatedAt:   workspace.CreatedAt,
			UpdatedAt:   workspace.UpdatedAt,
		}
	}

	return result, nil
}

func (s *WorkspaceServiceImpl) AddMemberToWorkspace(ctx context.Context,slug string,requsterId string,member dto.AddMemberDTO) (bool,error){

	ws,err := s.repo.GetSlug(ctx,slug)
	if err != nil{
		return false,fmt.Errorf("workspace not found")
	}

	is_authorised := false
	for _,value := range ws.Members{
		if(value.UserID == requsterId && (value.Role == schemas.RoleAdmin ||value.Role == schemas.RoleOwner)){
			is_authorised = true
			break
		}
	}

	if(!is_authorised){
		return false,fmt.Errorf("forbidden: only owners and admins can add members")
	}

	isalreadyMember,err:=s.repo.MemberExists(ctx,slug,member.UserID)
	if err != nil{
		return false,err
	}

	if(isalreadyMember){
		return  false,fmt.Errorf("user is already a member of this workspace")
	}

	memberPayload:=schemas.WorkspaceMember{
		UserID: member.UserID,
		AvatarURL: "/api/v1/users/" + member.UserID + "/avatar",
		Role: member.Role,
		JoinedAt: time.Now(),
	}
	
	ok,err := s.repo.AddMemberToWorkspace(ctx,slug,memberPayload)
	if err!= nil{
		return false,err
	}

	if(!ok){
		return false,fmt.Errorf("Error occured while adding member to the workspace, check with the db side")
	}

	return true,nil
}