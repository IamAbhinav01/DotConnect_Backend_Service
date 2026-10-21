package middlewares

import (
	"DotConnect/dto"
	"DotConnect/utils/json"
	"DotConnect/utils/response"
	"DotConnect/validator"
	"context"
	"net/http"
)

var PayloadContextKey contextKey = "payload"

func WorkSpaceCreation(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			var payload dto.WorkSpaceCreationDTO

			err := json.DecodeFromJSON(r, &payload)
			if err != nil {
				response.ErrorResponse(w, http.StatusBadRequest, "Error occured while decoding the json", err)
				return
			}
			validationErr := validator.Validate.Struct(payload)
			if validationErr != nil {
				response.ErrorResponse(w, http.StatusBadRequest, "Invalid Bad Request Payload", validationErr)
				return
			}

			reqContext := r.Context()
			ctx := context.WithValue(reqContext, PayloadContextKey, payload)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
}