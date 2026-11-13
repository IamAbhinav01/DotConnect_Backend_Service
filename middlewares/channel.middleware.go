package middlewares

import (
	"DotConnect/dto"
	"DotConnect/utils/json"
	"DotConnect/utils/response"
	"DotConnect/validator"
	"context"
	"net/http"
)

var ChannelContextKey contextKey = "Channelpayload"

func ChannelCreation(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			var payload dto.ChannelCreationDTO

			err := json.DecodeFromJSON(r, &payload)
			if err != nil {
				response.ErrorResponse(w, http.StatusBadRequest, "Error occured while decoding the json", err)
				return
			}

			t := true
			if payload.Public == nil{
				payload.Public = &t
			}

			validationErr := validator.Validate.Struct(payload)
			if validationErr != nil {
				response.ErrorResponse(w, http.StatusBadRequest, "Invalid Bad Request Payload", validationErr)
				return
			}

			reqContext := r.Context()
			ctx := context.WithValue(reqContext, ChannelContextKey, payload)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
}