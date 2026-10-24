package controllers

import (
	"DotConnect/dto"
	"DotConnect/middlewares"
	"DotConnect/services"
	"DotConnect/utils/response"
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


func (cntrl *WorkspaceController) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	
	dtoPayload, ok := r.Context().Value(middlewares.PayloadContextKey).(dto.WorkSpaceCreationDTO)
	if !ok {
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid payload in context", nil)
		return
	}

	// Get authenticated userID from context (set by AuthMiddleware)
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		response.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized: user ID missing", nil)
		return
	}

	workspace, err := cntrl.s.CreateWorkspace(r.Context(), userID, dtoPayload)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.SucessResponse(w, http.StatusCreated, "Workspace created successfully", workspace)
}


func (cntrl *WorkspaceController) GetAllUserWorkspace(w http.ResponseWriter, r *http.Request) {
	
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		response.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized: user ID missing", nil)
		return
	}
	
    workspaces, err := cntrl.s.GetUserWorkspace(r.Context(),userID)
    if err != nil {
        response.ErrorResponse(w, http.StatusInternalServerError, err.Error(), err)
        return
    }
    response.SucessResponse(w, http.StatusOK, "Workspaces fetched", workspaces)
}
