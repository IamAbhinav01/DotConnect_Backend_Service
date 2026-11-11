// DotConnect_Backend_Service/router/workspace.router.go
package router

import (
	"DotConnect/controllers"
	"DotConnect/middlewares"

	"github.com/go-chi/chi/v5"
)

type WorkspaceRouter struct {
	WorkspaceController *controllers.WorkspaceController
}

func NewWorkspaceRouter(ctrl *controllers.WorkspaceController) *WorkspaceRouter {
	return &WorkspaceRouter{
		WorkspaceController: ctrl,
	}
}

func (ws *WorkspaceRouter) Register(r chi.Router) {
	r.With(middlewares.AuthMiddleware, middlewares.WorkSpaceCreation).Post("/", ws.WorkspaceController.CreateWorkspace)
	r.With(middlewares.AuthMiddleware).Get("/",ws.WorkspaceController.GetAllUserWorkspace)
	r.With(middlewares.AuthMiddleware).Get("/{slug}/members",ws.WorkspaceController.GetAllMembers)
	r.With(middlewares.AuthMiddleware, middlewares.AddMemberValidation).Post("/{slug}/addMember", ws.WorkspaceController.AddMemberTOWorkspace)
	r.With(middlewares.AuthMiddleware,middlewares.AddMemberValidation).Post("/{slug}/members/promote",ws.WorkspaceController.PromoteMember)
	r.With(middlewares.AuthMiddleware,middlewares.AddMemberValidation).Post("/{slug}/members/demote",ws.WorkspaceController.DemoteMember)
	r.With(middlewares.AuthMiddleware,middlewares.InvitationCreation).Post("/{slug}/invites",ws.WorkspaceController.CreateInvite)
	r.With(middlewares.AuthMiddleware).Get("/invites/{code}",ws.WorkspaceController.InvitePreview)
	r.With(middlewares.AuthMiddleware).Get("/invite/{code}",ws.WorkspaceController.InvitePreview)
	r.With(middlewares.AuthMiddleware).Post("/invites/{code}/request",ws.WorkspaceController.JoinWorkspaceRequest)
	r.With(middlewares.AuthMiddleware).Post("/invite/{code}/request",ws.WorkspaceController.JoinWorkspaceRequest)
	r.With(middlewares.AuthMiddleware).Get("/{slug}/request",ws.WorkspaceController.GetAllPendingRequests)
	r.With(middlewares.AuthMiddleware).Post("/{slug}/request/{requestId}/approve",ws.WorkspaceController.ApproveJoinRequest)
	r.With(middlewares.AuthMiddleware).Post("/{slug}/request/{requestId}/reject",ws.WorkspaceController.RejectJoinRequest)
	r.With(middlewares.AuthMiddleware,middlewares.WorkSpaceUpdation).Post("/{slug}/update",ws.WorkspaceController.UpdateWorkspace)
	r.With(middlewares.AuthMiddleware, middlewares.RemoveMemberValidation).Delete("/{slug}/removeMember", ws.WorkspaceController.RemoveMemberFromWorkspace)
	r.With(middlewares.AuthMiddleware, middlewares.RemoveMemberValidation).Post("/{slug}/removeMember", ws.WorkspaceController.RemoveMemberFromWorkspace)
	r.With(middlewares.AuthMiddleware).Delete("/{slug}/removeWorkspace", ws.WorkspaceController.RemoveWorkspace)
	r.With(middlewares.AuthMiddleware).Post("/{slug}/removeWorkspace", ws.WorkspaceController.RemoveWorkspace)
}
