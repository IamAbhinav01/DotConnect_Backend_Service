package services

import (
	"DotConnect/dto"
	"DotConnect/repositories"
	"DotConnect/schemas"
	"context"
	"fmt"
	"strings"
	"time"
)

type WorkspaceService interface {
	CreateWorkspace(ctx context.Context,ownerId string,dto dto.WorkSpaceCreationDTO) (*schemas.WorkSpace, error)
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