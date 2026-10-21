package router

import "github.com/go-chi/chi/v5"

type WorkspaceRouter struct {
	// UserController controllers.UserController
	// RateLimiter    func(http.Handler) http.Handler
}

func (user *WorkspaceRouter) Register(r chi.Router) {

	// r.Get("/",)

	// r.With(user.RateLimiter, middleware.CreateUserRequestValidation).Post("/user/create", user.UserController.Create)
	// r.With(middleware.JWTMiddleware).Get("/users/{id}", user.UserController.GetUserByID)
	// r.With(user.RateLimiter).Get("/users", user.UserController.GetAllUsers)
	// r.Delete("/users/{id}", user.UserController.DeleteUserByID)
	// r.With(middleware.LoginUserRequestValidation).Post("/user/signin", user.UserController.Login)

	// r.Mount("/flights", middleware.FlightsProxy)
	// r.Mount("/booking", middleware.BookingsProxy)
	// r.Mount("/ai", middleware.AIProxy)
}

// func NewRouter(_userController *controllers.UserController,rateLimiter func(http.Handler) http.Handler) Router{
// 	return  &UserRouter{
// 		UserController: *_userController,
// 		RateLimiter:    rateLimiter,
// 	}
// }