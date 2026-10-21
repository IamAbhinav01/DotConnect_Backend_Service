package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)
type AppRouters struct {
	WorkspaceRouter *WorkspaceRouter
}

func SetupRouter(routers AppRouters) *chi.Mux {

	chiRouter := chi.NewRouter()
	
	chiRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "x-access-token", "x-idempotency-key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))


	workSpaceRouter := chi.NewRouter()
	routers.WorkspaceRouter.Register(workSpaceRouter)

	chiRouter.Mount("/api/v1/workspace", workSpaceRouter)
	chiRouter.Mount("/workspace", workSpaceRouter)
	chiRouter.Mount("/", workSpaceRouter)

	
	return chiRouter

}