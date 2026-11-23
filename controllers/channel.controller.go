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

type ChannelController struct {
	serv services.ChannelService
}

func NewChannelController(_s services.ChannelService)* ChannelController{ 
	return &ChannelController{
		serv: _s,
	}
}


func (cntrl *ChannelController) CreateChannel(w http.ResponseWriter, r *http.Request) {
	

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.ErrorResponse(w, http.StatusBadRequest, "Workspace slug is required", nil)
		return
	}

	dtoPayload, ok := r.Context().Value(middlewares.ChannelContextKey).(dto.ChannelCreationDTO)
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

	channel, err := cntrl.serv.CreateChannel(r.Context(), userID,slug ,dtoPayload)
	if err != nil {
		if strings.HasPrefix(err.Error(), "forbidden") {
			response.ErrorResponse(w, http.StatusForbidden, err.Error(), err)
			return
    	}
		response.ErrorResponse(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	response.SucessResponse(w, http.StatusCreated, "Channel created successfully", channel)
}

func (cntrl *ChannelController) GetAllChannels(w http.ResponseWriter, r *http.Request){

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.ErrorResponse(w, http.StatusBadRequest, "Workspace slug is required", nil)
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		response.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized: user ID missing", nil)
		return
	}

	channels , err := cntrl.serv.GetAllChannels(r.Context(),userID,slug)
	if err != nil {
		if strings.HasPrefix(err.Error(), "forbidden") {
			response.ErrorResponse(w, http.StatusForbidden, err.Error(), err)
			return
    	}
		response.ErrorResponse(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	response.SucessResponse(w, http.StatusCreated, "Channels fetched successfully", channels)
}