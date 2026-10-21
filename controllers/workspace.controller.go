package controllers

import (
	"DotConnect/dto"
	"DotConnect/middlewares"
	"DotConnect/services"
	"net/http"
)

type WorkspaceController struct {
	s services.WorkspaceService
}

func NewWorkspaceController(_s services.WorkspaceService)*WorkspaceController{
	return &WorkspaceController{
		s: _s,
	}
}

func (cntrl *WorkspaceController)CreateWorkspace(w http.ResponseWriter, r *http.Request){

	dtoPayload := r.Context().Value(middlewares.PayloadContextKey).(dto.WorkSpaceCreationDTO)

	dtoPayload.

}