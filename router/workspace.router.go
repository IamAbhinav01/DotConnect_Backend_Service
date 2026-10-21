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
}
