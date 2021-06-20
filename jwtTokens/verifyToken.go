package jwtTokens

import (
	"fmt"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

func Verify_Token(request_token string) (*jwt.Token, error) {
	tokenString := request_token
	err1 := godotenv.Load()
	if err1 != nil {
		log.Fatal("error loading .env file")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESSKEY")), nil //enter secret key
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
