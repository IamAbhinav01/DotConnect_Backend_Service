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
)

type WorkspaceService interface {
	CreateWorkspace(ctx context.Context,ownerId string,dto dto.WorkSpaceCreationDTO) (*schemas.WorkSpace, error)
	GetUserWorkspace(ctx context.Context,userId string)([]dto.WorkSpaceResponseDTO, error)
}

type WorkspaceServiceImpl struct {
	repo repositories.WorkspaceRepository
}

func NewWorkspaceService(_repo repositories.WorkspaceRepository) *WorkspaceServiceImpl {
	return &WorkspaceServiceImpl{
		repo: _repo,
	}
}

func(s *WorkspaceServiceImpl) CreateWorkspace(ctx context.Context,ownerId string,dto dto.WorkSpaceCreationDTO) (*schemas.WorkSpace, error){
	
	name := strings.TrimSpace(dto.Name)
	if name == ""{
		return nil,fmt.Errorf("workspace name cannot be empty")
	}
	description := dto.Description
	slug := dto.Slug
	
	if strings.TrimSpace(slug) == ""{
		slug = generator.Slugify(dto.Name)
	}


	ws := &schemas.WorkSpace{
		Name: name,
		Description: description,
		Slug: slug,
	}

	ws.Members = []schemas.WorkspaceMember{
		{
			UserID:   ownerId,
			Role:     schemas.RoleOwner,
			JoinedAt: time.Now(),
		},
	}
	

	_,err:=s.repo.CreateWorkspace(ctx,ws)
	if err!= nil{
		return nil,err
	}
	return ws,nil
}
func (s *WorkspaceServiceImpl) GetUserWorkspace(ctx context.Context,userId string)([]dto.WorkSpaceResponseDTO, error){
	workspaces,err := s.repo.GetAllWorkspaces(ctx,userId)
	if err != nil{
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