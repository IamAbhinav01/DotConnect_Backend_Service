package router

import (
	"DotConnect/controllers"
	"DotConnect/middlewares"

	"github.com/go-chi/chi/v5"
)

type ChannelRouter struct {
	ChannelController *controllers.ChannelController
}

func NewChannelRouter(ctrl *controllers.ChannelController) *ChannelRouter {
	return &ChannelRouter{
		ChannelController: ctrl,
	}
}


func (ws *ChannelRouter) Register(r chi.Router) {
	r.With(middlewares.AuthMiddleware, middlewares.ChannelCreation).Post("/", ws.ChannelController.CreateChannel)
}