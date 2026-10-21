package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Router interface {
	Register(chi.Router)
}
func SetupRouter(Router Router) *chi.Mux {

	chiRouter := chi.NewRouter()
	
	chiRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "x-access-token", "x-idempotency-key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))


	apiRouter := chi.NewRouter()
	Router.Register(apiRouter)

	chiRouter.Mount("/api/v1/workspace", apiRouter)
	chiRouter.Mount("/api/v1/channels",apiRouter)
	chiRouter.Mount("/api/v1/messages",apiRouter)

	
	return chiRouter

}