package auth

import (
	"DotConnect/config"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET = config.GetString("JWT_SECRET")

func VerifyJWTToken(token string) (*jwt.Token, error){
	jwt_token,err := jwt.Parse(token,func(token *jwt.Token) (interface{}, error){
		return JWT_SECRET,nil
	})

	if err != nil{
		return  nil,fmt.Errorf("Error occured while jwt token")
	}

	if !jwt_token.Valid{
		return nil,fmt.Errorf("Invalid token recieved")
	}

	return jwt_token,nil
}