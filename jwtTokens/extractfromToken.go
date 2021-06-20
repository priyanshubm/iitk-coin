package jwtTokens

import (
	"github.com/dgrijalva/jwt-go"
)

func Extract_data(user_token string) (string, error) {
	token, err := Verify_Token(user_token)
	if err != nil {
		return " ", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		roll_no, _ := claims["user_roll_no"].(string)
		return roll_no, err
	}

	return " ", err

}
