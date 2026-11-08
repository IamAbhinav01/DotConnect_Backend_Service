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
	invite, err := cntrl.s.CreateInvite(r.Context(),slug,userID,dtoPayload)
	if err != nil {
		if strings.HasPrefix(err.Error(), "forbidden") {
			response.ErrorResponse(w, http.StatusForbidden, err.Error(), err)
			return
		}
		response.ErrorResponse(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	
	response.SucessResponse(w, http.StatusOK, "Invite link created", invite)
}

func(cntrl *WorkspaceController) InvitePreview(w http.ResponseWriter, r *http.Request){
	code := chi.URLParam(r, "code")
	if code == "" {
		response.ErrorResponse(w, http.StatusBadRequest, "Workspace invite code is required", nil)
		return
	}

	// Get authenticated userID from context (set by AuthMiddleware)
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		response.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized: user ID missing", nil)
		return
	}
	preview,err := cntrl.s.FetchInvitation(r.Context(),code,userID)
	if err != nil {
		if strings.HasPrefix(err.Error(), "forbidden") {
			response.ErrorResponse(w, http.StatusForbidden, err.Error(), err)
			return
		}
		response.ErrorResponse(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	
	response.SucessResponse(w, http.StatusOK, "Invite preview fetched", preview)
}

func (cntrl *WorkspaceController) JoinWorkspaceRequest(w http.ResponseWriter, r *http.Request){

	code := chi.URLParam(r, "code")
	if code == "" {
		response.ErrorResponse(w, http.StatusBadRequest, "Workspace invitation code is required", nil)
		return
	}

	
	// Get authenticated userID from context (set by AuthMiddleware)
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		response.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized: user ID missing", nil)
		return
	}
	err := cntrl.s.CreateJoinRequest(r.Context(),code,userID)
	if err != nil {
		if strings.HasPrefix(err.Error(), "forbidden") {
			response.ErrorResponse(w, http.StatusForbidden, err.Error(), err)
			return
		}
		response.ErrorResponse(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	
	response.SucessResponse(w, http.StatusOK, "Join request submitted. Waiting for owner approval.", nil)
}

func (cntrl *WorkspaceController) GetAllPendingRequests(w http.ResponseWriter, r *http.Request) {
	
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.ErrorResponse(w, http.StatusBadRequest, "Workspace code is required", nil)
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		response.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized: user ID missing", nil)
		return
	}
	
    requests, err := cntrl.s.GetAllRequests(r.Context(),slug,userID)
    if err != nil {
	
        response.ErrorResponse(w, http.StatusInternalServerError, err.Error(), err)
        return
    }
    response.SucessResponse(w, http.StatusOK, "Pending requests fetched", requests)
}

func (cntrl *WorkspaceController) ApproveJoinRequest(w http.ResponseWriter,r *http.Request){
	slug := chi.URLParam(r, "slug")
	reqId := chi.URLParam(r,"requestId")
	if slug == "" {
		response.ErrorResponse(w, http.StatusBadRequest, "Workspace code is required", nil)
		return
	}
	if reqId == "" {
		response.ErrorResponse(w, http.StatusBadRequest, "request id is required", nil)
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		response.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized: user ID missing", nil)
		return
	}
	request,err := cntrl.s.ApproveJoinRequest(r.Context(),slug,userID,reqId)
	if err != nil{
		 response.ErrorResponse(w, http.StatusInternalServerError, err.Error(), err)
        return
	}
	response.SucessResponse(w, http.StatusOK, "Join request approved", request)
}

func (cntrl *WorkspaceController) RejectJoinRequest(w http.ResponseWriter,r *http.Request){
	slug := chi.URLParam(r, "slug")
	reqId := chi.URLParam(r,"requestId")
	if slug == "" {
		response.ErrorResponse(w, http.StatusBadRequest, "Workspace code is required", nil)
		return
	}
	if reqId == "" {
		response.ErrorResponse(w, http.StatusBadRequest, "request id is required", nil)
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		response.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized: user ID missing", nil)
		return
	}
	msg,err := cntrl.s.RejectJoinRequest(r.Context(),slug,userID,reqId)
	if err != nil{
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	response.SucessResponse(w, http.StatusOK, "Join request rejected", msg)
}