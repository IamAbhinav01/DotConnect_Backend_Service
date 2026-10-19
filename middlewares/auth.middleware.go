package middlewares

import (
	"DotConnect/auth"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string
const userIDKey contextKey = "userID"

func Auth_Handler(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split( r.Header.Get("Authorization"),"Bearer ")

		if len(authHeader) != 2{
			fmt.Println("Malformed Header")
			http.Error(w,"Invalid Header recieved",http.StatusBadRequest)
			return
		}else{
			token:=authHeader[1]
			jwt_token,err:=auth.VerifyJWTToken(token)
			
			if err != nil{
				fmt.Println("Invalid Token recieved")
				http.Error(w,"Invalid Token recieved",http.StatusUnauthorized)
				return
			}

			claims,ok:=jwt_token.Claims.(jwt.MapClaims)

			if !ok{
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userID, ok := claims["userId"].(string)
			if !ok || strings.TrimSpace(userID) == "" {
				http.Error(w, "token has no userId",http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))

		}

	})
}