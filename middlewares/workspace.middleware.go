package middlewares

import (
	"DotConnect/dto"
	"DotConnect/utils/json"
	"DotConnect/utils/response"
	"DotConnect/validator"
	"context"
	"io"
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

var AddToMemberContextKey contextKey = "addmemberContextKey"

func AddMemberValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			var payload dto.AddMemberDTO

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
			ctx := context.WithValue(reqContext, AddToMemberContextKey, payload)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
}

var CreateInvitationContextKey = "createInvitationContextKey"

func InvitationCreation(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			var payload dto.CreateInviteDTO

			// All invite fields are optional — an empty body is valid.
			// Decode only when there's actually a body; treat io.EOF as "no body".
			err := json.DecodeFromJSON(r, &payload)
			if err != nil && err != io.EOF {
				response.ErrorResponse(w, http.StatusBadRequest, "Error occured while decoding the json", err)
				return
			}

			// Only validate struct when something was actually provided
			if err == nil {
				validationErr := validator.Validate.Struct(payload)
				if validationErr != nil {
					response.ErrorResponse(w, http.StatusBadRequest, "Invalid Bad Request Payload", validationErr)
					return
				}
			}

			reqContext := r.Context()
			ctx := context.WithValue(reqContext, CreateInvitationContextKey, payload)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
}