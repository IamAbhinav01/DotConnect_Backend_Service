package controllers

import (
	"DotConnect/dto"
	"DotConnect/middlewares"
	"DotConnect/services"
	"DotConnect/utils/response"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
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

func (cntrl *WorkspaceController) GetAllMembers(w http.ResponseWriter, r *http.Request) {
	
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.ErrorResponse(w, http.StatusBadRequest, "Workspace slug is required", nil)
		return
	}
	
    workspaces, err := cntrl.s.GetAllMembers(r.Context(),slug)
    if err != nil {
        response.ErrorResponse(w, http.StatusInternalServerError, err.Error(), err)
        return
    }
    response.SucessResponse(w, http.StatusOK, "members fetched", workspaces)
}


func (cntrl *WorkspaceController) AddMemberTOWorkspace(w http.ResponseWriter, r *http.Request){

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.ErrorResponse(w, http.StatusBadRequest, "Workspace slug is required", nil)
		return
	}

	dtoPayload, ok := r.Context().Value(middlewares.AddToMemberContextKey).(dto.AddMemberDTO)
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

	_, err := cntrl.s.AddMemberToWorkspace(r.Context(), slug,userID, dtoPayload)
	if err != nil {
		if err.Error() == "forbidden: only owners and admins can add members" {
			response.ErrorResponse(w, http.StatusForbidden, err.Error(), err)
			return
		}
		response.ErrorResponse(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	response.SucessResponse(w, http.StatusCreated, "Member added to workspace successfully", nil)
}

func (cntrl *WorkspaceController) PromoteMember(w http.ResponseWriter, r *http.Request){

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.ErrorResponse(w, http.StatusBadRequest, "Workspace slug is required", nil)
		return
	}

	dtoPayload, ok := r.Context().Value(middlewares.AddToMemberContextKey).(dto.AddMemberDTO)
	
	if !ok {
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid payload in context", nil)
		return
	}

	targetUser := dtoPayload.UserID

	// Get authenticated userID from context (set by AuthMiddleware)
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		response.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized: user ID missing", nil)
		return
	}

	_, err := cntrl.s.PromoteMember(r.Context(), slug,userID,targetUser )
	if err != nil {
		if err.Error() == "forbidden: only owners and admins can promote members" {
			response.ErrorResponse(w, http.StatusForbidden, err.Error(), err)
			return
		}
		response.ErrorResponse(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	response.SucessResponse(w, http.StatusOK, "Member promoted to admin successfully", nil)
}

func (cntrl *WorkspaceController) DemoteMember(w http.ResponseWriter, r *http.Request){

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.ErrorResponse(w, http.StatusBadRequest, "Workspace slug is required", nil)
		return
	}

	dtoPayload, ok := r.Context().Value(middlewares.AddToMemberContextKey).(dto.AddMemberDTO)
	
	if !ok {
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid payload in context", nil)
		return
	}

	targetUser := dtoPayload.UserID

	// Get authenticated userID from context (set by AuthMiddleware)
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		response.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized: user ID missing", nil)
		return
	}

	_, err := cntrl.s.DemoteMember(r.Context(), slug,userID,targetUser )
	if err != nil {
		if err.Error() == "forbidden: only owners can demote admins" {
			response.ErrorResponse(w, http.StatusForbidden, err.Error(), err)
			return
		}
		response.ErrorResponse(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	response.SucessResponse(w, http.StatusOK, "admin demoted successfully", nil)
}

func (cntrl *WorkspaceController) CreateInvite(w http.ResponseWriter, r *http.Request){

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.ErrorResponse(w, http.StatusBadRequest, "Workspace slug is required", nil)
		return
	}

	dtoPayload, ok := r.Context().Value(middlewares.CreateInvitationContextKey).(dto.CreateInviteDTO)
	
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
	err := cntrl.s.CreateInvite(r.Context(),slug,userID,dtoPayload)
	if err != nil {
		if strings.HasPrefix(err.Error(), "forbidden") {
			response.ErrorResponse(w, http.StatusForbidden, err.Error(), err)
			return
		}
		response.ErrorResponse(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	
	response.SucessResponse(w, http.StatusOK, "Invite link created", nil)
}