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
}
