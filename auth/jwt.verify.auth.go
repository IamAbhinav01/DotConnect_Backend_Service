package auth

import (
	"DotConnect/config"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var internalSecret = []byte(config.GetString("INTERNAL_SECRET"))

type Claims struct {
	UserID string `json:"userId"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func VerifyIdentityToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return internalSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid or expired identity token")
	}

	if claims.Issuer != "api-gateway" || claims.Audience[0] != "dotconnect-go-service" {
		return nil, fmt.Errorf("token audience/issuer mismatch")
	}

	return claims, nil
}
